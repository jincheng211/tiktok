package ffmpeg

import (
	"bytes"
	"os"
	"os/exec"
)

func GetVideoCover(video []byte) ([]byte, error) {

	// 保存视频数据到临时文件
	tmpfile, err := os.CreateTemp("", "video_*.mp4")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	_, err = tmpfile.Write(video)
	if err != nil {
		return nil, err
	}

	// 使用FFmpeg命令行工具获取视频封面
	cmd := exec.Command("ffmpeg", "-i", tmpfile.Name(), "-vframes", "1", "-f", "image2", "-")
	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
