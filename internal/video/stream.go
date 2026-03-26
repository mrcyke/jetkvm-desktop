package video

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/pion/rtp/codecs"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media"
	"github.com/pion/webrtc/v4/pkg/media/h264reader"
	"github.com/pion/webrtc/v4/pkg/media/samplebuilder"
)

type Frame struct {
	Image image.Image
	At    time.Time
}

type Stream struct {
	mu         sync.RWMutex
	latest     *Frame
	frameCh    chan Frame
	closeOnce  sync.Once
	cancelFunc context.CancelFunc
}

func NewStream() *Stream {
	return &Stream{
		frameCh: make(chan Frame, 4),
	}
}

func (s *Stream) Latest() *Frame {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.latest
}

func (s *Stream) publish(frame Frame) {
	s.mu.Lock()
	s.latest = &frame
	s.mu.Unlock()

	select {
	case s.frameCh <- frame:
	default:
	}
}

func (s *Stream) Frames() <-chan Frame {
	return s.frameCh
}

func (s *Stream) Close() {
	s.closeOnce.Do(func() {
		if s.cancelFunc != nil {
			s.cancelFunc()
		}
		close(s.frameCh)
	})
}

func AttachRemoteTrack(parent context.Context, track *webrtc.TrackRemote) (*Stream, error) {
	if track.Codec().MimeType != webrtc.MimeTypeH264 {
		return nil, fmt.Errorf("unsupported codec %s", track.Codec().MimeType)
	}

	ctx, cancel := context.WithCancel(parent)
	stream := NewStream()
	stream.cancelFunc = cancel

	ffmpegCmd := exec.CommandContext(ctx,
		"ffmpeg",
		"-loglevel", "error",
		"-fflags", "nobuffer",
		"-flags", "low_delay",
		"-f", "h264",
		"-i", "pipe:0",
		"-f", "image2pipe",
		"-vcodec", "png",
		"pipe:1",
	)

	stdin, err := ffmpegCmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := ffmpegCmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := ffmpegCmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		defer stream.Close()
		defer stdin.Close()

		sb := samplebuilder.New(32, &codecs.H264Packet{}, track.Codec().ClockRate)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			pkt, _, readErr := track.ReadRTP()
			if readErr != nil {
				return
			}
			sb.Push(pkt)
			for {
				sample := sb.Pop()
				if sample == nil {
					break
				}
				payload := ensureAnnexB(sample.Data)
				if len(payload) == 0 {
					continue
				}
				if _, err := stdin.Write(payload); err != nil {
					return
				}
			}
		}
	}()

	go func() {
		defer stream.Close()
		reader := bufio.NewReader(stdout)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			img, err := png.Decode(reader)
			if err != nil {
				if err == io.EOF || ctx.Err() != nil {
					return
				}
				return
			}
			stream.publish(Frame{Image: img, At: time.Now()})
		}
	}()

	return stream, nil
}

func StartTestPattern(ctx context.Context, width, height, fps int, track *webrtc.TrackLocalStaticSample) error {
	cmd := exec.CommandContext(ctx,
		"ffmpeg",
		"-loglevel", "error",
		"-f", "lavfi",
		"-i", fmt.Sprintf("testsrc=size=%dx%d:rate=%d", width, height, fps),
		"-pix_fmt", "yuv420p",
		"-an",
		"-c:v", "libx264",
		"-preset", "ultrafast",
		"-tune", "zerolatency",
		"-x264-params", "repeat-headers=1:aud=1:keyint=30:min-keyint=30:scenecut=0",
		"-f", "h264",
		"pipe:1",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	reader, err := h264reader.NewReader(stdout)
	if err != nil {
		return err
	}

	frameDuration := time.Second / time.Duration(fps)
	go func() {
		defer cmd.Wait()

		var au bytes.Buffer
		flush := func() bool {
			if au.Len() == 0 {
				return true
			}
			err := track.WriteSample(media.Sample{
				Data:     append([]byte(nil), au.Bytes()...),
				Duration: frameDuration,
			})
			au.Reset()
			return err == nil
		}

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			nal, readErr := reader.NextNAL()
			if readErr != nil {
				return
			}

			if nal.UnitType == h264reader.NalUnitTypeAUD && au.Len() > 0 {
				if !flush() {
					return
				}
			}
			au.Write([]byte{0x00, 0x00, 0x00, 0x01})
			au.Write(nal.Data)
		}
	}()

	return nil
}

func ensureAnnexB(data []byte) []byte {
	if len(data) == 0 {
		return nil
	}
	if bytes.HasPrefix(data, []byte{0x00, 0x00, 0x00, 0x01}) || bytes.HasPrefix(data, []byte{0x00, 0x00, 0x01}) {
		return data
	}
	return append([]byte{0x00, 0x00, 0x00, 0x01}, data...)
}
