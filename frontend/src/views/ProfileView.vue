<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useUserStore } from '@/stores/user';
import { socialAPI, videoAPI, likeAPI } from '@/api';
import type { Profile, Video } from '@/types';
import { ArrowLeft, Settings, Heart, MessageCircle, Grid, Play, LogOut, Edit3 } from 'lucide-vue-next';

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();

const profile = ref<Profile | null>(null);
const isLoading = ref(true);
const error = ref<string | null>(null);
const activeTab = ref<'works' | 'likes'>('works');
const videos = ref<Video[]>([]);
const likedVideos = ref<Video[]>([]);

const isCurrentUser = computed(() => {
  if (!profile.value) return false;
  return profile.value.account_id === userStore.accountId;
});

const formatNumber = (num: number) => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'w';
  }
  return num.toString();
};

const loadProfile = async () => {
  isLoading.value = true;
  error.value = null;
  try {
    const accountId = parseInt(route.params.id as string);
    if (isNaN(accountId) || accountId <= 0) {
      throw new Error('无效的用户ID');
    }
    profile.value = await socialAPI.getProfile(accountId);
    await loadVideos(accountId);
  } catch (err) {
    console.error('Failed to load profile:', err);
    error.value = err instanceof Error ? err.message : '加载失败';
  } finally {
    isLoading.value = false;
  }
};

const loadVideos = async (accountId: number) => {
  try {
    const response = await videoAPI.getUserVideos(accountId, 1, 20);
    videos.value = response.videos || [];
  } catch (err) {
    console.error('Failed to load videos:', err);
    videos.value = [];
  }
};

const loadLikedVideos = async () => {
  if (!userStore.isLoggedIn()) {
    likedVideos.value = [];
    return;
  }
  
  try {
    const response = await likeAPI.getLikedVideos(0, 20);
    likedVideos.value = response.videos || [];
  } catch (err) {
    console.error('Failed to load liked videos:', err);
    likedVideos.value = [];
  }
};

const switchTab = async (tab: 'works' | 'likes') => {
  activeTab.value = tab;
  if (tab === 'likes' && likedVideos.value.length === 0) {
    await loadLikedVideos();
  }
};

const playVideo = (videoId: number) => {
  router.push({ path: '/', query: { video_id: videoId } });
};

const handleFollow = async () => {
  if (!profile.value || !userStore.isLoggedIn()) {
    router.push('/login');
    return;
  }
  try {
    if (profile.value.is_followed) {
      await socialAPI.unfollow(profile.value.account_id);
      profile.value.is_followed = false;
    } else {
      await socialAPI.follow(profile.value.account_id);
      profile.value.is_followed = true;
    }
  } catch (err) {
    console.error('Follow error:', err);
  }
};

const handleLogout = async () => {
  await userStore.logout();
  router.push('/login');
};

const handleMessage = () => {
  if (!profile.value || !userStore.isLoggedIn()) {
    router.push('/login');
    return;
  }
  router.push({ path: '/message', query: { user_id: profile.value.account_id } });
};

onMounted(() => {
  loadProfile();
});
</script>

<template>
  <div class="profile-page">
    <header class="profile-header">
      <button class="back-btn" @click="router.back()">
        <ArrowLeft :size="24" />
      </button>
      <h1 class="profile-title">个人主页</h1>
      <button class="settings-btn">
        <Settings :size="24" />
      </button>
    </header>

    <div v-if="isLoading" class="loading-state">
      <div class="spinner"></div>
    </div>

    <div v-else-if="error" class="error-state">
      <p>{{ error }}</p>
      <button class="retry-btn" @click="loadProfile">重试</button>
    </div>

    <div v-else-if="profile" class="profile-content">
      <div class="profile-info">
        <div class="avatar-section">
          <div class="avatar-container">
            <img 
              v-if="profile.avatar_url" 
              :src="profile.avatar_url" 
              class="avatar" 
              :alt="profile.username"
            />
            <div v-else class="avatar-placeholder">{{ profile.username.charAt(0) }}</div>
          </div>
          <div class="stats-row">
            <div class="stat-item">
              <span class="stat-value">{{ videos.length }}</span>
              <span class="stat-label">作品</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ formatNumber(profile.follower_count) }}</span>
              <span class="stat-label">粉丝</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ formatNumber(profile.following_count) }}</span>
              <span class="stat-label">关注</span>
            </div>
          </div>
        </div>

        <div class="user-info">
          <h2 class="username">{{ profile.username }}</h2>
          <p class="bio">{{ profile.bio || '暂无简介' }}</p>
        </div>

        <div class="action-buttons">
          <button 
            v-if="!isCurrentUser"
            class="action-btn"
            :class="{ following: profile.is_followed }"
            @click="handleFollow"
          >
            {{ profile.is_followed ? '已关注' : '+ 关注' }}
          </button>
          <button 
            v-if="!isCurrentUser"
            class="action-btn secondary"
            @click="handleMessage"
          >
            发私信
          </button>
          <button 
            v-if="isCurrentUser"
            class="action-btn"
            @click="router.push('/profile/edit')"
          >
            <Edit3 :size="16" />
            <span>编辑资料</span>
          </button>
          <button 
            v-if="isCurrentUser"
            class="action-btn secondary"
            @click="handleLogout"
          >
            <LogOut :size="16" />
            <span>退出登录</span>
          </button>
        </div>

        <div class="tabs">
          <button 
            class="tab" 
            :class="{ active: activeTab === 'works' }"
            @click="switchTab('works')"
          >
            <Grid :size="20" />
            <span>作品</span>
          </button>
          <button 
            class="tab" 
            :class="{ active: activeTab === 'likes' }"
            @click="switchTab('likes')"
          >
            <Play :size="20" />
            <span>喜欢</span>
          </button>
        </div>
      </div>

      <div v-if="activeTab === 'works'" class="video-grid">
        <template v-if="videos.length > 0">
          <div 
            v-for="video in videos" 
            :key="video.id" 
            class="video-item"
            @click="playVideo(video.id)"
          >
            <div class="video-thumbnail">
              <img :src="video.cover_url" :alt="video.title" />
              <div class="play-overlay">
                <Play :size="40" />
              </div>
            </div>
            <div class="video-stats">
              <span class="stat"><Heart :size="14" /> {{ formatNumber(video.like_count) }}</span>
              <span class="stat"><MessageCircle :size="14" /> {{ formatNumber(video.comment_count) }}</span>
            </div>
          </div>
        </template>
        <div v-else class="empty-state">
          <Grid :size="64" />
          <p>暂无作品</p>
        </div>
      </div>

      <div v-else class="video-grid">
        <template v-if="likedVideos.length > 0">
          <div 
            v-for="video in likedVideos" 
            :key="video.id" 
            class="video-item"
            @click="playVideo(video.id)"
          >
            <div class="video-thumbnail">
              <img :src="video.cover_url" :alt="video.title" />
              <div class="play-overlay">
                <Play :size="40" />
              </div>
            </div>
            <div class="video-stats">
              <span class="stat"><Heart :size="14" /> {{ formatNumber(video.like_count) }}</span>
              <span class="stat"><MessageCircle :size="14" /> {{ formatNumber(video.comment_count) }}</span>
            </div>
          </div>
        </template>
        <div v-else class="empty-state">
          <Heart :size="64" />
          <p>暂无喜欢</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.profile-page {
  min-height: 100vh;
  background: #000;
}

.profile-header {
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

.back-btn, .settings-btn {
  background: none;
  border: none;
  color: #fff;
  cursor: pointer;
  padding: 8px;
}

.profile-title {
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

.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  padding: 20px;
}

.error-state p {
  color: #999;
  font-size: 16px;
  margin-bottom: 16px;
}

.retry-btn {
  background: #ff2d55;
  border: none;
  border-radius: 20px;
  color: #fff;
  padding: 10px 24px;
  font-size: 14px;
  cursor: pointer;
}

.profile-content {
  padding-top: 56px;
}

.profile-info {
  padding: 20px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.avatar-section {
  display: flex;
  align-items: center;
  gap: 32px;
  margin-bottom: 20px;
}

.avatar-container {
  width: 100px;
  height: 100px;
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
  font-size: 40px;
  font-weight: bold;
}

.stats-row {
  flex: 1;
  display: flex;
  justify-content: space-around;
}

.stat-item {
  text-align: center;
}

.stat-value {
  display: block;
  font-size: 20px;
  font-weight: bold;
}

.stat-label {
  font-size: 12px;
  color: #999;
}

.user-info {
  margin-bottom: 16px;
}

.username {
  font-size: 20px;
  font-weight: bold;
  margin: 0 0 8px 0;
}

.bio {
  font-size: 14px;
  color: #999;
  margin: 0;
}

.action-buttons {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}

.action-btn {
  flex: 1;
  height: 36px;
  background: #ff2d55;
  border: none;
  border-radius: 8px;
  color: #fff;
  font-size: 14px;
  font-weight: bold;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.action-btn.secondary {
  background: rgba(255, 255, 255, 0.1);
}

.action-btn.following {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

.tabs {
  display: flex;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.tab {
  flex: 1;
  height: 48px;
  background: none;
  border: none;
  color: #999;
  font-size: 14px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  border-bottom: 2px solid transparent;
}

.tab.active {
  color: #fff;
  border-bottom-color: #fff;
}

.video-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 2px;
}

.video-item {
  aspect-ratio: 1;
  position: relative;
  cursor: pointer;
}

.video-thumbnail {
  width: 100%;
  height: 100%;
  position: relative;
}

.video-thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.play-overlay {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: #fff;
  opacity: 0;
  transition: opacity 0.3s;
}

.video-item:hover .play-overlay {
  opacity: 1;
}

.video-stats {
  position: absolute;
  bottom: 4px;
  left: 4px;
  right: 4px;
  display: flex;
  justify-content: space-between;
  padding: 4px 8px;
  background: rgba(0, 0, 0, 0.6);
}

.stat {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #fff;
}

.empty-state {
  grid-column: 1 / -1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #999;
}

.empty-state p {
  margin-top: 16px;
  font-size: 16px;
}
</style>