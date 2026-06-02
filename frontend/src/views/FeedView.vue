<script setup lang="ts">import { ref, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/stores/user';
import { feedAPI, likeAPI, commentAPI, socialAPI, videoAPI } from '@/api';
import type { FeedItem } from '@/types';
import { Heart, MessageCircle, Send, Home, Search, PlusCircle, MessageCircle as MsgIcon, User, Flame, Users, Hash } from 'lucide-vue-next';
const router = useRouter();
const userStore = useUserStore();
const feedItems = ref<FeedItem[]>([]);
const currentIndex = ref(0);
const isLoading = ref(false);
const feedType = ref<'hot' | 'latest' | 'following'>('hot');
const showComments = ref(false);
const newComment = ref('');
const currentVideoId = ref(0);
const comments = ref<any[]>([]);
const loadFeed = async (type: 'hot' | 'latest' | 'following') => {
 feedType.value = type;
 isLoading.value = true;
 try {
 let response;
 switch (type) {
 case 'hot':
 response = await feedAPI.getHot(10);
 break;
 case 'latest':
 response = await feedAPI.getLatest(0, 10);
 break;
 case 'following':
 response = await feedAPI.getFollowing(0, 10);
 break;
 }
 feedItems.value = response.items;
 currentIndex.value = 0;
 }
 catch (err) {
 console.error('Failed to load feed:', err);
 }
 finally {
 isLoading.value = false;
 }
};
const handleLike = async (videoId: number) => {
 if (!userStore.isLoggedIn()) {
 router.push('/login');
 return;
 }
 try {
 const item = feedItems.value.find(i => i.video_id === videoId);
 if (item) {
 if (item.is_liked) {
 const response = await likeAPI.unlike(videoId);
 item.like_count = response.like_count;
 item.is_liked = false;
 }
 else {
 const response = await likeAPI.like(videoId);
 item.like_count = response.like_count;
 item.is_liked = true;
 }
 }
 }
 catch (err) {
 console.error('Like error:', err);
 }
};
const handleFollow = async (accountId: number) => {
 if (!userStore.isLoggedIn()) {
 router.push('/login');
 return;
 }
 try {
 const item = feedItems.value.find(i => i.account_id === accountId);
 if (!item) return;
 
 if (item.is_following) {
 await socialAPI.unfollow(accountId);
 item.is_following = false;
 } else {
 await socialAPI.follow(accountId);
 item.is_following = true;
 }
 }
 catch (err) {
 console.error('Follow error:', err);
 }
};
const openComments = async (videoId: number) => {
 currentVideoId.value = videoId;
 showComments.value = true;
 await loadComments(videoId);
};
const closeComments = () => {
 showComments.value = false;
 newComment.value = '';
 comments.value = [];
};
const loadComments = async (videoId: number) => {
 try {
 const response = await commentAPI.list(videoId, 1, 20);
 comments.value = response.comments;
 } catch (err) {
 console.error('Load comments error:', err);
 }
};
const submitComment = async () => {
 if (!newComment.value.trim() || !userStore.isLoggedIn())
 return;
 try {
 await commentAPI.create(currentVideoId.value, newComment.value);
 newComment.value = '';
 await loadComments(currentVideoId.value);
 const item = feedItems.value.find(i => i.video_id === currentVideoId.value);
 if (item) {
 item.comment_count++;
 }
 }
 catch (err) {
 console.error('Comment error:', err);
 }
};
const scrollToVideo = (direction: 'up' | 'down') => {
 const container = document.querySelector('.feed-container');
 if (!container)
 return;
 const videoHeight = window.innerHeight;
 if (direction === 'down' && currentIndex.value < feedItems.value.length - 1) {
 currentIndex.value++;
 container.scrollBy({ top: videoHeight, behavior: 'smooth' });
 }
 else if (direction === 'up' && currentIndex.value > 0) {
 currentIndex.value--;
 container.scrollBy({ top: -videoHeight, behavior: 'smooth' });
 }
};

const handleVideoPlay = async (videoId: number) => {
 try {
 await videoAPI.reportView(videoId);
 const item = feedItems.value.find(i => i.video_id === videoId);
 if (item) {
 item.view_count++;
 }
 }
 catch (err) {
 console.error('Report view error:', err);
 }
};
const handleScroll = () => {
 const container = document.querySelector('.feed-container');
 if (!container)
 return;
 const scrollTop = container.scrollTop;
 const videoHeight = window.innerHeight;
 currentIndex.value = Math.round(scrollTop / videoHeight);
};
const formatNumber = (num: number) => {
 if (num >= 10000) {
 return (num / 10000).toFixed(1) + 'w';
 }
 return num.toString();
};
const formatTime = (timeStr: string) => {
 if (!timeStr) return '刚刚';
 const date = new Date(timeStr);
 const now = new Date();
 const diff = now.getTime() - date.getTime();
 const minutes = Math.floor(diff / 60000);
 const hours = Math.floor(diff / 3600000);
 const days = Math.floor(diff / 86400000);
 if (minutes < 1) return '刚刚';
 if (minutes < 60) return minutes + '分钟前';
 if (hours < 24) return hours + '小时前';
 if (days < 7) return days + '天前';
 return date.toLocaleDateString('zh-CN');
};
onMounted(() => {
 loadFeed('hot');
 const container = document.querySelector('.feed-container');
 container?.addEventListener('scroll', handleScroll);
});
onUnmounted(() => {
 const container = document.querySelector('.feed-container');
 container?.removeEventListener('scroll', handleScroll);
});
</script>

<template>
  <div class="feed-page">
    <header class="feed-header">
      <div class="header-left">
        <button class="nav-btn" @click="loadFeed('hot')" :class="{ active: feedType === 'hot' }">
          <Flame :size="20" />
          <span>热门</span>
        </button>
        <button class="nav-btn" @click="loadFeed('latest')" :class="{ active: feedType === 'latest' }">
          <Hash :size="20" />
          <span>最新</span>
        </button>
        <button 
          v-if="userStore.isLoggedIn()" 
          class="nav-btn" 
          @click="loadFeed('following')" 
          :class="{ active: feedType === 'following' }"
        >
          <Users :size="20" />
          <span>关注</span>
        </button>
      </div>
      <div class="header-center">
        <h1 class="logo">🎵 抖音</h1>
      </div>
      <div class="header-right">
        <button class="icon-btn" @click="router.push('/search')">
          <Search :size="24" />
        </button>
      </div>
    </header>

    <div v-if="isLoading" class="loading-state">
      <div class="spinner"></div>
      <p>加载中...</p>
    </div>

    <div v-else class="feed-container" @wheel.prevent="(e) => scrollToVideo(e.deltaY > 0 ? 'down' : 'up')">
      <div v-for="item in feedItems" :key="item.video_id" class="video-card">
        <div class="video-wrapper">
          <video 
            :src="item.video_url" 
            class="video-player"
            autoplay
            loop
            muted
            playsinline
            :poster="item.cover_url"
            @play="handleVideoPlay(item.video_id)"
          ></video>
          <div class="video-overlay">
            <div class="video-info">
              <h3>{{ item.title }}</h3>
              <p class="description">{{ item.description }}</p>
              <div class="tags">
                <span v-for="tag in item.tags.split(',')" :key="tag" class="tag">#{{ tag }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="video-sidebar">
          <button 
            class="sidebar-btn avatar-btn" 
            @click="router.push(`/profile/${item.account_id}`)"
          >
            <img 
              v-if="item.avatar_url" 
              :src="item.avatar_url" 
              class="avatar" 
              :alt="item.username"
            />
            <div v-else class="avatar-placeholder">{{ item.username.charAt(0) }}</div>
          </button>
          <p class="username">{{ item.username }}</p>
          <button 
            class="sidebar-btn follow-btn" 
            :class="{ following: item.is_following }"
            @click="handleFollow(item.account_id)"
          >
            {{ item.is_following ? '已关注' : '+ 关注' }}
          </button>

          <div class="actions">
            <button class="sidebar-btn action-btn" @click="handleLike(item.video_id)">
              <Heart :size="28" :class="{ liked: item.is_liked }" />
              <span>{{ formatNumber(item.like_count) }}</span>
            </button>

            <button class="sidebar-btn action-btn" @click="openComments(item.video_id)">
              <MessageCircle :size="28" />
              <span>{{ formatNumber(item.comment_count) }}</span>
            </button>

            </div>
        </div>

        <div class="video-footer">
          <div class="view-count">
            <span>{{ formatNumber(item.view_count) }} 播放</span>
          </div>
        </div>
      </div>
    </div>

    <nav class="bottom-nav">
      <button class="nav-item active" @click="router.push('/')">
        <Home :size="24" />
        <span>首页</span>
      </button>
      <button class="nav-item" @click="router.push('/search')">
        <Search :size="24" />
        <span>搜索</span>
      </button>
      <button class="nav-item upload-btn" @click="router.push('/upload')">
        <PlusCircle :size="24" />
      </button>
      <button class="nav-item" @click="router.push('/message')">
        <MsgIcon :size="24" />
        <span>消息</span>
      </button>
      <button class="nav-item" @click="userStore.isLoggedIn() ? router.push(`/profile/${userStore.accountId}`) : router.push('/login')">
            <User :size="24" />
            <span>我</span>
          </button>
    </nav>

    <div v-if="showComments" class="comments-modal" @click.self="closeComments">
      <div class="comments-container">
        <div class="comments-header">
          <h3>评论</h3>
          <button class="close-btn" @click="closeComments">✕</button>
        </div>
        <div class="comments-list">
          <div 
            v-for="comment in comments" 
            :key="comment.id" 
            class="comment-item"
          >
            <div class="comment-avatar">{{ comment.username?.charAt(0) || '?' }}</div>
            <div class="comment-content">
              <div class="comment-header">
                <span class="comment-username">{{ comment.username }}</span>
                <span class="comment-time">{{ formatTime(comment.created_at) }}</span>
              </div>
              <p>{{ comment.content }}</p>
            </div>
          </div>
          <div v-if="comments.length === 0" class="no-comments">
            <p>暂无评论，快来抢沙发吧~</p>
          </div>
        </div>
        <div class="comments-input">
          <input 
            v-model="newComment" 
            type="text" 
            placeholder="说点什么..." 
            class="comment-input"
            @keyup.enter="submitComment"
          />
          <button class="send-btn" @click="submitComment">
            <Send :size="20" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.feed-page {
  min-height: 100vh;
  background: #000;
  display: flex;
  flex-direction: column;
}

.feed-header {
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

.header-left {
  display: flex;
  gap: 16px;
}

.nav-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  background: none;
  border: none;
  color: #fff;
  font-size: 14px;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 20px;
  transition: background 0.3s;
}

.nav-btn:hover, .nav-btn.active {
  background: rgba(255, 255, 255, 0.1);
}

.header-center .logo {
  font-size: 24px;
  font-weight: bold;
  margin: 0;
}

.icon-btn {
  background: none;
  border: none;
  color: #fff;
  cursor: pointer;
  padding: 8px;
}

.feed-container {
  margin-top: 56px;
  margin-bottom: 60px;
  overflow-y: auto;
  scroll-snap-type: y mandatory;
}

.video-card {
  height: calc(100vh - 116px);
  scroll-snap-align: start;
  position: relative;
}

.video-wrapper {
  width: 100%;
  height: 100%;
  position: relative;
}

.video-player {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.video-overlay {
  position: absolute;
  bottom: 120px;
  left: 16px;
  right: 80px;
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.8));
  padding: 20px 16px;
}

.video-info {
  color: #fff;
}

.video-info h3 {
  font-size: 18px;
  margin: 0 0 8px 0;
}

.description {
  font-size: 14px;
  margin: 0 0 12px 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag {
  font-size: 12px;
  color: #00f5ff;
}

.video-sidebar {
  position: absolute;
  right: 16px;
  bottom: 120px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.sidebar-btn {
  background: none;
  border: none;
  color: #fff;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.avatar-btn {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  overflow: hidden;
  border: 2px solid #ff2d55;
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
  font-size: 20px;
  font-weight: bold;
}

.username {
  font-size: 12px;
  margin: 0;
}

.follow-btn {
  background: #ff2d55;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: bold;
}

.actions {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.action-btn {
  gap: 8px;
}

.action-btn span {
  font-size: 12px;
}

.action-btn .liked {
  color: #ff2d55;
}

.video-footer {
  position: absolute;
  bottom: 16px;
  left: 16px;
  color: rgba(255, 255, 255, 0.7);
  font-size: 12px;
}

.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: rgba(0, 0, 0, 0.95);
  display: flex;
  align-items: center;
  justify-content: space-around;
  padding: 0 16px;
  z-index: 100;
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  background: none;
  border: none;
  color: #999;
  cursor: pointer;
  font-size: 10px;
}

.nav-item.active {
  color: #ff2d55;
}

.upload-btn {
  width: 44px;
  height: 44px;
  background: linear-gradient(135deg, #ff0050, #ff2d55);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: -20px;
  box-shadow: 0 4px 12px rgba(255, 45, 85, 0.4);
}

.loading-state {
  display: flex;
  flex-direction: column;
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

.comments-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: flex-end;
  z-index: 200;
}

.comments-container {
  width: 100%;
  max-height: 70vh;
  background: #1a1a1a;
  border-radius: 24px 24px 0 0;
  display: flex;
  flex-direction: column;
}

.comments-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.comments-header h3 {
  margin: 0;
}

.close-btn {
  background: none;
  border: none;
  color: #999;
  font-size: 20px;
  cursor: pointer;
}

.comments-list {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.comment-item {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}

.comment-avatar {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #ff0050, #ff2d55);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
}

.comment-content {
  flex: 1;
}

.comment-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.comment-username {
  font-weight: bold;
  font-size: 14px;
}

.comment-time {
  font-size: 12px;
  color: #666;
}

.comments-input {
  display: flex;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.comment-input {
  flex: 1;
  height: 40px;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 20px;
  padding: 0 16px;
  color: #fff;
  font-size: 14px;
}

.send-btn {
  width: 40px;
  height: 40px;
  background: #ff2d55;
  border: none;
  border-radius: 50%;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

@media (min-width: 768px) {
  .feed-container {
    width: 480px;
    margin-left: auto;
    margin-right: auto;
  }
}
</style>