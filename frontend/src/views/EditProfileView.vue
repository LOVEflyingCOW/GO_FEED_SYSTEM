<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/stores/user';
import { accountAPI, socialAPI } from '@/api';
import type { Profile } from '@/types';
import { ArrowLeft, Camera, Save } from 'lucide-vue-next';

const router = useRouter();
const userStore = useUserStore();

const profile = ref<Profile | null>(null);
const username = ref('');
const bio = ref('');
const avatarFile = ref<File | null>(null);
const avatarPreview = ref('');
const isLoading = ref(false);
const isSubmitting = ref(false);

const loadProfile = async () => {
  isLoading.value = true;
  try {
    profile.value = await socialAPI.getProfile(userStore.accountId);
    username.value = profile.value.username;
    bio.value = profile.value.bio || '';
  } catch (err) {
    console.error('Failed to load profile:', err);
  } finally {
    isLoading.value = false;
  }
};

const handleAvatarChange = (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (target.files && target.files.length > 0) {
    avatarFile.value = target.files[0];
    // 添加即时预览
    const reader = new FileReader();
    reader.onload = (e) => {
      avatarPreview.value = e.target?.result as string;
    };
    reader.readAsDataURL(target.files[0]);
  }
};

const handleSubmit = async () => {
  if (!username.value.trim()) {
    alert('请输入用户名');
    return;
  }

  isSubmitting.value = true;
  try {
    if (avatarFile.value) {
      const formData = new FormData();
      formData.append('avatar', avatarFile.value);
      
      console.log('准备上传头像，文件名:', avatarFile.value.name);
      console.log('文件大小:', avatarFile.value.size, 'bytes');
      console.log('文件类型:', avatarFile.value.type);
      
      const uploadPromise = accountAPI.uploadAvatar(formData);
      const timeoutPromise = new Promise((_, reject) => 
        setTimeout(() => reject(new Error('上传超时')), 30000)
      );
      
      const result = await Promise.race([uploadPromise, timeoutPromise]);
      console.log('头像上传成功:', result);
    }

    const hasUsernameChanged = username.value !== profile.value?.username;
    const hasBioChanged = bio.value !== profile.value?.bio;

    if (hasUsernameChanged || hasBioChanged) {
      await accountAPI.updateProfile(hasUsernameChanged ? username.value : undefined, bio.value || undefined);
      if (hasUsernameChanged) {
        userStore.username = username.value;
        localStorage.setItem('username', username.value);
      }
    }

    alert('资料修改成功');
    router.back();
  } catch (err: any) {
    console.error('头像上传失败:', err);
    console.error('错误响应:', err.response);
    console.error('错误状态:', err.response?.status);
    console.error('错误数据:', err.response?.data);
    const errorMsg = err.response?.data?.message || err.message || '修改失败';
    alert(errorMsg);
  } finally {
    isSubmitting.value = false;
  }
};

onMounted(() => {
  if (!userStore.isLoggedIn()) {
    router.push('/login');
    return;
  }
  loadProfile();
});
</script>

<template>
  <div class="edit-profile-page">
    <header class="edit-profile-header">
      <button class="back-btn" @click="router.back()">
        <ArrowLeft :size="24" />
      </button>
      <h1 class="page-title">编辑资料</h1>
      <button class="save-btn" @click="handleSubmit" :disabled="isSubmitting">
        <Save :size="20" />
      </button>
    </header>

    <div v-if="isLoading" class="loading-state">
      <div class="spinner"></div>
    </div>

    <div v-else class="edit-profile-content">
      <div class="form-section">
        <div class="avatar-upload-section">
          <label class="avatar-upload-label">
            <div class="avatar-preview">
              <img 
                v-if="avatarPreview || profile?.avatar_url" 
                :src="avatarPreview || profile?.avatar_url" 
                class="avatar"
              />
              <div v-else class="avatar-placeholder">
                {{ username.charAt(0) }}
              </div>
              <div class="upload-overlay">
                <Camera :size="24" />
              </div>
            </div>
            <input 
              type="file" 
              accept="image/*" 
              class="avatar-input"
              @change="handleAvatarChange"
            />
          </label>
          <p class="avatar-hint">点击更换头像</p>
        </div>

        <div class="form-group">
          <label class="form-label">用户名</label>
          <input 
            v-model="username" 
            type="text" 
            class="form-input"
            placeholder="请输入用户名"
            maxlength="20"
          />
        </div>

        <div class="form-group">
          <label class="form-label">简介</label>
          <textarea 
            v-model="bio" 
            class="form-textarea"
            placeholder="介绍一下你自己..."
            maxlength="150"
          ></textarea>
          <span class="char-count">{{ bio.length }}/150</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.edit-profile-page {
  min-height: 100vh;
  background: #000;
}

.edit-profile-header {
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

.back-btn, .save-btn {
  background: none;
  border: none;
  color: #fff;
  cursor: pointer;
  padding: 8px;
}

.save-btn:disabled {
  opacity: 0.5;
}

.page-title {
  margin: 0;
  font-size: 18px;
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
}

.spinner {
  width: 48px;
  height: 48px;
  border: 4px solid rgba(255, 255, 255, 0.1);
  border-top-color: #ff2d55;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.edit-profile-content {
  padding: 72px 16px 20px;
}

.form-section {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  padding: 20px;
}

.avatar-upload-section {
  text-align: center;
  margin-bottom: 24px;
}

.avatar-upload-label {
  cursor: pointer;
}

.avatar-preview {
  position: relative;
  width: 120px;
  height: 120px;
  margin: 0 auto;
  border-radius: 50%;
  overflow: hidden;
  border: 3px solid #ff2d55;
}

.avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #ff0050, #ff2d55);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48px;
  font-weight: bold;
}

.upload-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s;
}

.avatar-upload-label:hover .upload-overlay {
  opacity: 1;
}

.avatar-input {
  display: none;
}

.avatar-hint {
  color: #999;
  font-size: 14px;
  margin: 12px 0 0 0;
}

.form-group {
  margin-bottom: 20px;
}

.form-label {
  display: block;
  color: #999;
  font-size: 14px;
  margin-bottom: 8px;
}

.form-input {
  width: 100%;
  height: 44px;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 8px;
  padding: 0 16px;
  color: #fff;
  font-size: 16px;
  outline: none;
}

.form-textarea {
  width: 100%;
  height: 120px;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 8px;
  padding: 12px 16px;
  color: #fff;
  font-size: 16px;
  outline: none;
  resize: none;
}

.char-count {
  display: block;
  text-align: right;
  color: #666;
  font-size: 12px;
  margin-top: 8px;
}
</style>