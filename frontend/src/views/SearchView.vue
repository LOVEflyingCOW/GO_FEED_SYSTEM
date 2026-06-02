<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { feedAPI, accountAPI } from '@/api';
import { Search, X, TrendingUp, User, Video } from 'lucide-vue-next';

const router = useRouter();
const searchQuery = ref('');
const searchHistory = ref(['美食', '旅行', '搞笑', '音乐']);
const trendingTags = [
  { tag: '热门挑战', count: '12.5w' },
  { tag: '每日一笑', count: '8.2w' },
  { tag: '美食探店', count: '6.8w' },
  { tag: '旅行日记', count: '5.3w' },
];
const searchResults = ref<any[]>([]);
const usersResults = ref<any[]>([]);
const isSearching = ref(false);
const activeTab = ref<'video' | 'user'>('video');

const handleSearch = async () => {
  if (!searchQuery.value.trim()) return;
  
  isSearching.value = true;
  try {
    const [videoRes, userRes] = await Promise.all([
      feedAPI.search(searchQuery.value),
      accountAPI.searchUsers(searchQuery.value, 10)
    ]);
    
    const videoData = videoRes;
    const userData = userRes;
    
    searchResults.value = videoData.items || [];
    usersResults.value = userData || [];
  } catch (err) {
    console.error('Search error:', err);
    searchResults.value = [];
    usersResults.value = [];
  } finally {
    isSearching.value = false;
  }
};

const clearSearch = () => {
  searchQuery.value = '';
  searchResults.value = [];
  usersResults.value = [];
};

const selectTag = (tag: string) => {
  searchQuery.value = tag;
  handleSearch();
};

const formatNumber = (num: number) => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'w';
  }
  return num.toString();
};
</script>

<template>
  <div class="search-page">
    <header class="search-header">
      <div class="search-box">
        <Search :size="20" />
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="搜索视频或用户" 
          class="search-input"
          @keyup.enter="handleSearch"
        />
        <button v-if="searchQuery" class="clear-btn" @click="clearSearch">
          <X :size="18" />
        </button>
        <button v-if="searchQuery" class="search-btn" @click="handleSearch">
          搜索
        </button>
      </div>
    </header>

    <div class="search-content">
      <div v-if="!searchQuery" class="search-home">
        <div class="trending-section">
          <div class="section-header">
            <TrendingUp :size="18" />
            <h3>热门搜索</h3>
          </div>
          <div class="trending-list">
            <div 
              v-for="(item, index) in trendingTags" 
              :key="item.tag"
              class="trending-item"
              @click="selectTag(item.tag)"
            >
              <span class="trending-rank">{{ index + 1 }}</span>
              <span class="trending-tag">{{ item.tag }}</span>
              <span class="trending-count">{{ item.count }}</span>
            </div>
          </div>
        </div>

        <div class="history-section">
          <div class="section-header">
            <h3>搜索历史</h3>
            <button class="clear-history">清空</button>
          </div>
          <div class="history-tags">
            <span 
              v-for="tag in searchHistory" 
              :key="tag" 
              class="history-tag"
              @click="selectTag(tag)"
            >
              {{ tag }}
            </span>
          </div>
        </div>
      </div>

      <div v-else class="search-results">
        <div class="tabs">
          <button 
            class="tab" 
            :class="{ active: activeTab === 'video' }"
            @click="activeTab = 'video'"
          >
            <Video :size="18" />
            <span>视频</span>
          </button>
          <button 
            class="tab" 
            :class="{ active: activeTab === 'user' }"
            @click="activeTab = 'user'"
          >
            <User :size="18" />
            <span>用户</span>
          </button>
        </div>

        <div v-if="isSearching" class="loading-state">
          <div class="spinner"></div>
          <p>搜索中...</p>
        </div>

        <div v-else-if="activeTab === 'video'" class="results-list video-results">
          <div 
            v-for="result in searchResults" 
            :key="result.video_id"
            class="video-result-item"
            @click="router.push(`/profile/${result.account_id}`)"
          >
            <div class="video-thumbnail">
              <img :src="result.cover_url" alt="视频封面" />
              <div class="video-duration">{{ result.duration }}</div>
            </div>
            <div class="video-info">
              <h4>{{ result.title }}</h4>
              <p class="video-desc">{{ result.description }}</p>
              <div class="video-stats">
                <span class="stat"><Video :size="14" /> {{ formatNumber(result.view_count) }}</span>
                <span class="stat"><span>💬</span> {{ formatNumber(result.comment_count) }}</span>
              </div>
            </div>
          </div>
          <div v-if="searchResults.length === 0" class="no-results">
            <p>没有找到相关视频</p>
          </div>
        </div>

        <div v-else class="results-list user-results">
          <div 
            v-for="user in usersResults" 
            :key="user.id"
            class="user-result-item"
            @click="router.push(`/profile/${user.id}`)"
          >
            <div class="user-avatar">{{ user.username?.charAt(0) || '?' }}</div>
            <div class="user-info">
              <span class="user-name">{{ user.username }}</span>
            </div>
          </div>
          <div v-if="usersResults.length === 0" class="no-results">
            <p>没有找到相关用户</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.search-page {
  min-height: 100vh;
  background: #000;
}

.search-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  background: #000;
  padding: 12px 16px;
  z-index: 100;
}

.search-box {
  display: flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 24px;
  padding: 8px 16px;
}

.search-box svg {
  color: #999;
  margin-right: 12px;
}

.search-input {
  flex: 1;
  background: none;
  border: none;
  color: #fff;
  font-size: 16px;
  outline: none;
}

.clear-btn {
  background: none;
  border: none;
  color: #999;
  cursor: pointer;
  padding: 4px;
}

.search-btn {
  background: #ff2d55;
  border: none;
  border-radius: 16px;
  color: #fff;
  padding: 6px 16px;
  font-size: 14px;
  cursor: pointer;
  margin-left: 8px;
}

.search-content {
  padding-top: 60px;
}

.search-home {
  padding: 16px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}

.section-header h3 {
  margin: 0;
  font-size: 16px;
  color: #fff;
}

.clear-history {
  margin-left: auto;
  background: none;
  border: none;
  color: #999;
  font-size: 14px;
  cursor: pointer;
}

.trending-list {
  margin-bottom: 24px;
}

.trending-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  cursor: pointer;
}

.trending-rank {
  width: 24px;
  height: 24px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: bold;
}

.trending-item:nth-child(1) .trending-rank {
  background: linear-gradient(135deg, #ff0050, #ff2d55);
}

.trending-item:nth-child(2) .trending-rank {
  background: linear-gradient(135deg, #ff7a00, #ff9500);
}

.trending-item:nth-child(3) .trending-rank {
  background: linear-gradient(135deg, #ffb800, #ffcc00);
}

.trending-tag {
  flex: 1;
  font-size: 16px;
}

.trending-count {
  color: #999;
  font-size: 14px;
}

.history-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.history-tag {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 6px 12px;
  font-size: 14px;
  cursor: pointer;
}

.search-results {
  padding: 16px;
}

.tabs {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}

.tab {
  display: flex;
  align-items: center;
  gap: 8px;
  background: none;
  border: none;
  color: #999;
  font-size: 16px;
  cursor: pointer;
  padding: 8px 0;
  border-bottom: 2px solid transparent;
}

.tab.active {
  color: #fff;
  border-bottom-color: #ff2d55;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid rgba(255, 255, 255, 0.1);
  border-top-color: #ff2d55;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loading-state p {
  color: #999;
  margin-top: 12px;
}

.results-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.video-result-item {
  display: flex;
  gap: 12px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  overflow: hidden;
}

.video-thumbnail {
  position: relative;
  width: 120px;
  height: 80px;
  flex-shrink: 0;
}

.video-thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.video-duration {
  position: absolute;
  bottom: 4px;
  right: 4px;
  background: rgba(0, 0, 0, 0.8);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
}

.video-info {
  flex: 1;
  padding: 8px;
}

.video-info h4 {
  margin: 0 0 4px 0;
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.video-desc {
  margin: 0 0 8px 0;
  color: #999;
  font-size: 12px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.video-stats {
  display: flex;
  gap: 16px;
}

.stat {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #999;
  font-size: 12px;
}

.user-result-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  cursor: pointer;
}

.user-avatar {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #ff0050, #ff2d55);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: bold;
}

.user-name {
  font-size: 16px;
}

.no-results {
  text-align: center;
  padding: 40px;
}

.no-results p {
  color: #999;
}
</style>