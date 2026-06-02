<script setup lang="ts">import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { videoAPI } from '@/api';
import { ArrowLeft, Upload, Image, X, Check } from 'lucide-vue-next';
const router = useRouter();
const videoFile = ref<File | null>(null);
const coverFile = ref<File | null>(null);
const title = ref('');
const description = ref('');
const tags = ref('');
const isUploading = ref(false);
const progress = ref(0);
const videoPreview = ref('');
const coverPreview = ref('');
const videoInput = ref<HTMLInputElement | null>(null);
const coverInput = ref<HTMLInputElement | null>(null);

const handleVideoUpload = (event: Event) => {
 const target = event.target as HTMLInputElement;
 const file = target.files?.[0];
 if (file) {
 videoFile.value = file;
 videoPreview.value = URL.createObjectURL(file);
 }
};
const handleCoverUpload = (event: Event) => {
 const target = event.target as HTMLInputElement;
 const file = target.files?.[0];
 if (file) {
 coverFile.value = file;
 coverPreview.value = URL.createObjectURL(file);
 }
};
const removeVideo = () => {
 if (videoPreview.value) {
 URL.revokeObjectURL(videoPreview.value);
 }
 videoFile.value = null;
 videoPreview.value = '';
};
const openVideoInput = () => {
 videoInput.value?.click();
};
const openCoverInput = () => {
 coverInput.value?.click();
};
const removeCover = () => {
 if (coverPreview.value) {
 URL.revokeObjectURL(coverPreview.value);
 }
 coverFile.value = null;
 coverPreview.value = '';
};
const handleSubmit = async () => {
 if (!videoFile.value || !title.value) {
 alert('请选择视频文件并填写标题');
 return;
 }
 isUploading.value = true;
 progress.value = 0;
 try {
 const formData = new FormData();
 formData.append('video', videoFile.value);
 if (coverFile.value) {
 formData.append('cover', coverFile.value);
 }
 formData.append('title', title.value);
 if (description.value) {
 formData.append('description', description.value);
 }
 if (tags.value) {
 formData.append('tags', tags.value);
 }
 const interval = setInterval(() => {
 progress.value += Math.random() * 20;
 if (progress.value > 90) {
 progress.value = 90;
 }
 }, 300);
 await videoAPI.upload(formData);
 clearInterval(interval);
 progress.value = 100;
 setTimeout(() => {
 router.push('/');
 }, 1000);
 }
 catch (err) {
 console.error('Upload failed:', err);
 alert('上传失败，请重试');
 }
 finally {
 isUploading.value = false;
 }
};
</script>

<template>
  <div class="upload-page">
    <header class="upload-header">
      <button class="back-btn" @click="router.back()">
        <ArrowLeft :size="24" />
      </button>
      <h1 class="upload-title">上传视频</h1>
      <button 
        class="upload-btn" 
        :disabled="!videoFile || !title || isUploading"
        @click="handleSubmit"
      >
        <Upload :size="18" />
        <span>发布</span>
      </button>
    </header>

    <div class="upload-content">
      <div class="video-upload-section">
        <h3>视频文件</h3>
        <div 
          v-if="!videoPreview" 
          class="upload-area"
          @click="openVideoInput"
        >
          <Upload :size="48" />
          <p>点击上传视频</p>
          <p class="hint">支持 MP4、MOV 格式</p>
          <input 
            ref="videoInput" 
            type="file" 
            accept="video/*" 
            class="hidden-input"
            @change="handleVideoUpload"
          />
        </div>
        <div v-else class="preview-container">
          <video :src="videoPreview" controls class="preview-video"></video>
          <button class="remove-btn" @click="removeVideo">
            <X :size="20" />
          </button>
        </div>
      </div>

      <div class="cover-upload-section">
        <h3>封面图片（可选）</h3>
        <div 
          v-if="!coverPreview" 
          class="upload-area"
          @click="openCoverInput"
        >
          <Image :size="48" />
          <p>点击上传封面</p>
          <p class="hint">支持 JPG、PNG 格式</p>
          <input 
            ref="coverInput" 
            type="file" 
            accept="image/*" 
            class="hidden-input"
            @change="handleCoverUpload"
          />
        </div>
        <div v-else class="preview-container">
          <img :src="coverPreview" class="preview-cover" />
          <button class="remove-btn" @click="removeCover">
            <X :size="20" />
          </button>
        </div>
      </div>

      <div class="form-section">
        <div class="form-group">
          <label>标题</label>
          <input 
            v-model="title" 
            type="text" 
            placeholder="添加视频标题..." 
            class="form-input"
            maxlength="100"
          />
        </div>

        <div class="form-group">
          <label>描述</label>
          <textarea 
            v-model="description" 
            placeholder="添加视频描述..." 
            class="form-textarea"
            maxlength="500"
          ></textarea>
        </div>

        <div class="form-group">
          <label>标签</label>
          <input 
            v-model="tags" 
            type="text" 
            placeholder="添加标签，用逗号分隔..." 
            class="form-input"
          />
          <p class="hint">如：美食,旅行,搞笑</p>
        </div>
      </div>
    </div>

    <div v-if="isUploading" class="upload-modal">
      <div class="upload-progress">
        <div class="progress-bar">
          <div class="progress-fill" :style="{ width: progress + '%' }"></div>
        </div>
        <p>{{ Math.round(progress) }}%</p>
        <p v-if="progress === 100" class="success-text">
          <Check :size="24" />
          上传成功！
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.upload-page {
  min-height: 100vh;
  background: #000;
}

.upload-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 56px;
  background: rgba(0, 0, 0, 0.9);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  z-index: 100;
}

.back-btn {
  background: none;
  border: none;
  color: #fff;
  cursor: pointer;
  padding: 8px;
}

.upload-title {
  margin: 0;
  font-size: 18px;
}

.upload-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #ff2d55;
  border: none;
  border-radius: 20px;
  color: #fff;
  padding: 8px 16px;
  font-size: 14px;
  font-weight: bold;
  cursor: pointer;
}

.upload-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.upload-content {
  padding: 72px 16px 16px;
}

.video-upload-section, .cover-upload-section {
  margin-bottom: 24px;
}

.video-upload-section h3, .cover-upload-section h3 {
  color: #999;
  font-size: 14px;
  margin: 0 0 12px 0;
}

.upload-area {
  border: 2px dashed rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  padding: 40px 20px;
  text-align: center;
  cursor: pointer;
  transition: border-color 0.3s;
}

.upload-area:hover {
  border-color: #ff2d55;
}

.upload-area p {
  margin: 8px 0 0 0;
  color: #999;
  font-size: 14px;
}

.hint {
  font-size: 12px !important;
  color: #666 !important;
}

.hidden-input {
  display: none;
}

.preview-container {
  position: relative;
  border-radius: 12px;
  overflow: hidden;
}

.preview-video {
  width: 100%;
  max-height: 400px;
  object-fit: contain;
  background: #000;
}

.preview-cover {
  width: 100%;
  max-height: 300px;
  object-fit: cover;
}

.remove-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 32px;
  height: 32px;
  background: rgba(0, 0, 0, 0.6);
  border: none;
  border-radius: 50%;
  color: #fff;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.form-section {
  margin-top: 16px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  color: #999;
  font-size: 14px;
  margin-bottom: 8px;
}

.form-input {
  width: 100%;
  height: 44px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  padding: 0 12px;
  color: #fff;
  font-size: 16px;
  outline: none;
}

.form-input:focus {
  border-color: #ff2d55;
}

.form-textarea {
  width: 100%;
  height: 100px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  padding: 12px;
  color: #fff;
  font-size: 16px;
  outline: none;
  resize: none;
}

.form-textarea:focus {
  border-color: #ff2d55;
}

.upload-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
}

.upload-progress {
  text-align: center;
  padding: 40px;
}

.progress-bar {
  width: 200px;
  height: 8px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 4px;
  overflow: hidden;
  margin: 0 auto 16px;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #ff0050, #ff2d55);
  transition: width 0.3s;
}

.upload-progress p {
  color: #fff;
  font-size: 18px;
  margin: 0;
}

.success-text {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: #00ff88 !important;
}
</style>