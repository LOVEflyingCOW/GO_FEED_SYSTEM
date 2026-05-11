package video

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// UploadService 视频上传服务
type UploadService struct {
	videoRepository *VideoRepository // 视频数据访问层
	uploadDir       string           // 文件上传目录
	baseURL         string           // 资源访问基础URL
}

// NewUploadService 创建上传服务实例

func NewUploadService(videoRepository *VideoRepository, uploadDir, baseURL string) *UploadService {
	return &UploadService{
		videoRepository: videoRepository,
		uploadDir:       uploadDir,
		baseURL:         baseURL,
	}
}

// ctx: 上下文对象
// coverFile: 封面文件（可选）
// 返回上传结果和错误信息
func (us *UploadService) UploadVideo(ctx context.Context, accountID uint, username, title, description, tags string, videoFile *multipart.FileHeader, coverFile *multipart.FileHeader) (*UploadVideoResponse, error) {
	// 确保上传目录存在
	if err := os.MkdirAll(us.uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// 保存视频文件
	videoPath, err := us.saveFile(videoFile, accountID, "video")
	if err != nil {
		return nil, fmt.Errorf("failed to save video: %w", err)
	}

	// 处理封面：优先使用用户上传的封面，否则从视频中提取
	var coverPath string
	if coverFile != nil {
		coverPath, err = us.saveFile(coverFile, accountID, "cover")
		if err != nil {
			return nil, fmt.Errorf("failed to save cover: %w", err)
		}
	} else {
		coverPath, err = us.extractCover(videoPath, accountID)
		if err != nil {
			return nil, fmt.Errorf("failed to extract cover: %w", err)
		}
	}

	// 获取视频时长
	duration, err := us.getVideoDuration(videoPath)
	if err != nil {
		log.Printf("warning: failed to get duration: %v", err)
		duration = 0
	}

	// 构建视频实体
	var coverURL string
	if coverPath != "" {
		coverURL = us.baseURL + "/" + coverPath
	}
	video := &Video{
		AccountID:   accountID,
		Username:    username,
		Title:       title,
		VideoPath:   videoPath,
		CoverPath:   coverPath,
		PlayURL:     us.baseURL + "/" + videoPath,
		CoverURL:    coverURL,
		Duration:    duration,
		Description: description,
		Tags:        tags,
	}

	// 保存视频记录到数据库
	if err := us.videoRepository.CreateVideo(ctx, video); err != nil {
		// 创建记录失败时清理已上传的文件
		os.Remove(filepath.Join(us.uploadDir, videoPath))
		if coverPath != "" {
			os.Remove(filepath.Join(us.uploadDir, coverPath))
		}
		return nil, fmt.Errorf("failed to create video record: %w", err)
	}

	// 返回上传结果
	return &UploadVideoResponse{
		VideoID:     video.ID,
		VideoURL:    us.baseURL + "/" + videoPath,
		CoverURL:    us.baseURL + "/" + coverPath,
		Title:       video.Title,
		Description: video.Description,
		Tags:        video.Tags,
	}, nil
}

// saveFile 保存上传的文件
// file: 上传的文件头
// fileType: 文件类型标识
// 返回保存后的文件名和错误信息
func (us *UploadService) saveFile(file *multipart.FileHeader, accountID uint, fileType string) (string, error) {
	// 生成唯一文件名：accountID_timestamp_pid.ext
	ext := filepath.Ext(file.Filename)
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%d_%d%s", accountID, timestamp, os.Getpid(), ext)

	filePath := filepath.Join(us.uploadDir, filename)

	// 打开源文件
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return filename, nil
}

// extractCover 从视频中提取封面
// videoPath: 视频文件路径（相对路径）
// accountID: 用户账号ID
// 返回封面文件名和错误信息
func (us *UploadService) extractCover(videoPath string, accountID uint) (string, error) {
	// 生成封面文件名
	coverFilename := fmt.Sprintf("%d_%d_cover.jpg", accountID, time.Now().Unix())
	coverPath := filepath.Join(us.uploadDir, coverFilename)

	// 使用ffmpeg提取第一帧作为封面
	err := ffmpeg.Input(filepath.Join(us.uploadDir, videoPath)).
		Output(coverPath, ffmpeg.KwArgs{"vframes": 1, "q:v": 2}).
		OverWriteOutput().
		Run()
	if err != nil {
		return "", err
	}

	return coverFilename, nil
}

// getVideoDuration 获取视频时长（秒）
// videoPath: 视频文件路径（相对路径）
// 返回视频时长（秒）和错误信息
func (us *UploadService) getVideoDuration(videoPath string) (int, error) {
	path := filepath.Join(us.uploadDir, videoPath)
	// 使用ffmpeg探测视频信息
	result, err := ffmpeg.Probe(path)
	if err != nil {
		return 0, err
	}

	resultStr := string(result) //方便调试
	var probeResult map[string]interface{}
	// 解析JSON格式的探测结果
	if err := json.Unmarshal([]byte(resultStr), &probeResult); err != nil {
		return 0, err
	}

	formatMap, ok := probeResult["format"].(map[string]interface{})
	if !ok {
		return 0, nil
	}

	durationVal := formatMap["duration"]
	if durationVal == nil {
		return 0, nil
	}

	// 将时长字符串转换为整数秒
	duration, err := strconv.ParseFloat(fmt.Sprintf("%v", durationVal), 64)
	if err != nil {
		return 0, err
	}

	return int(duration), nil
}

// DeleteVideoFiles 删除视频文件和封面文件
// videoPath: 视频文件路径（相对路径）
// coverPath: 封面文件路径（相对路径）
func (us *UploadService) DeleteVideoFiles(videoPath, coverPath string) error {
	if videoPath != "" {
		if err := os.Remove(filepath.Join(us.uploadDir, videoPath)); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	if coverPath != "" {
		if err := os.Remove(filepath.Join(us.uploadDir, coverPath)); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

// ExtractTags 从文本中提取标签（#开头的单词）
// text: 输入文本
// 返回逗号分隔的标签字符串
func ExtractTags(text string) string {
	var tags []string
	words := strings.Fields(text)
	for _, word := range words {
		if strings.HasPrefix(word, "#") && len(word) > 1 {
			tags = append(tags, strings.TrimPrefix(word, "#"))
		}
	}
	return strings.Join(tags, ",")
}
