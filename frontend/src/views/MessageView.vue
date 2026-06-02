<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useUserStore } from '@/stores/user';
import { messageAPI, accountAPI } from '@/api';
import { ArrowLeft, Send, MoreVertical, User, Search, X } from 'lucide-vue-next';

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();

const conversations = ref<any[]>([]);
const messages = ref<any[]>([]);
const currentConversation = ref<number | null>(null);
const newMessage = ref('');
const showSearch = ref(false);
const searchQuery = ref('');
const searchResults = ref<any[]>([]);
const isLoading = ref(false);

const loadConversations = async () => {
  try {
    const response = await messageAPI.getConversations();
    conversations.value = response.map((conv: any) => ({
      id: conv.user_id,
      name: conv.username,
      avatar: conv.avatar_url,
      last_message: conv.last_message,
      unread: conv.unread_count,
      time: conv.updated_at ? new Date(conv.updated_at).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' }) : '刚刚'
    }));
  } catch (err) {
    console.error('Failed to load conversations:', err);
  }
};

const initSSE = () => {
  if (!userStore.isLoggedIn()) return;
  
  const token = localStorage.getItem('token');
  if (!token) return;
  
  const connectSSE = async () => {
    abortController = new AbortController();
    
    try {
      const response = await fetch('/api/v1/sse', {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Accept': 'text/event-stream'
        },
        credentials: 'include',
        signal: abortController.signal
      });
      
      if (!response.ok) {
        console.error('SSE connection failed:', response.status);
        return;
      }
      
      const reader = response.body?.getReader();
      if (!reader) {
        console.error('Failed to get reader');
        return;
      }
      
      const decoder = new TextDecoder('utf-8');
      let buffer = '';
      
      console.log('SSE connected');
      
      while (true) {
        const { done, value } = await reader.read();
        if (done) {
          console.log('SSE connection closed');
          setTimeout(connectSSE, 3000);
          break;
        }
        
        buffer += decoder.decode(value, { stream: true });
        const messages = buffer.split('\n\n');
        buffer = messages.pop() || '';
        
        for (const message of messages) {
          if (message.includes('data:')) {
            try {
              const lines = message.split('\n');
              let jsonStr = '';
              for (const line of lines) {
                if (line.startsWith('data:')) {
                  jsonStr = line.substring(5).trim();
                  break;
                }
              }
              if (jsonStr) {
                console.log('Received SSE message:', jsonStr);
                const data = JSON.parse(jsonStr);
                if (data.type === 'message') {
                  handleNewMessage(data);
                }
              }
            } catch (err) {
              console.error('Failed to parse SSE message:', err);
            }
          }
        }
      }
    } catch (err) {
      console.error('SSE error:', err);
      setTimeout(connectSSE, 3000);
    }
  };
  
  connectSSE();
};

const handleNewMessage = (notification: any) => {
  const fromId = notification.from_id;
  const content = notification.content;
  
  // 更新消息列表
  if (currentConversation.value === fromId) {
    messages.value.push({
      id: Date.now(),
      sender_id: fromId,
      content: content,
      time: new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' }),
      isMe: false
    });
  }
  
  // 更新会话列表
  const existingConv = conversations.value.find(c => c.id === fromId);
  if (existingConv) {
    existingConv.last_message = content;
    existingConv.unread++;
    existingConv.time = '刚刚';
  } else {
    // 如果是新会话，添加到列表
    conversations.value.unshift({
      id: fromId,
      name: notification.from || 'Unknown',
      avatar: '',
      last_message: content,
      unread: 1,
      time: '刚刚'
    });
  }
};

let abortController: AbortController | null = null;

const closeSSE = () => {
  if (abortController) {
    abortController.abort();
    abortController = null;
    console.log('SSE disconnected');
  }
};

const openChatWithUser = async (userId: number) => {
  try {
    const userData = await accountAPI.getAccount(userId);
    const userName = userData.username;
    
    const existingConv = conversations.value.find(c => c.id === userId);
    if (existingConv) {
      selectConversation(userId);
    } else {
      conversations.value.unshift({
        id: userId,
        name: userName,
        avatar: userData.avatar_url || '',
        lastMessage: '',
        time: '刚刚',
        unread: 0
      });
      selectConversation(userId);
    }
  } catch (err) {
    console.error('Failed to get user info:', err);
  }
};

const loadMessages = async (otherId: number) => {
  try {
    const response = await messageAPI.getMessages(otherId);
    messages.value = response.messages;
  } catch (err) {
    console.error('Failed to load messages:', err);
  }
};

const selectConversation = (id: number) => {
  currentConversation.value = id;
  loadMessages(id);
};

const sendMessage = async () => {
  if (!newMessage.value.trim() || !currentConversation.value) return;
  
  try {
    await messageAPI.send(currentConversation.value, newMessage.value);
    messages.value.push({
      id: Date.now(),
      sender_id: userStore.accountId,
      content: newMessage.value,
      time: new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' }),
      isMe: true
    });
    newMessage.value = '';
  } catch (err) {
    console.error('Send message error:', err);
  }
};

const searchUsers = async () => {
    if (!searchQuery.value.trim()) {
      return;
    }
    
    isLoading.value = true;
    try {
      const response = await accountAPI.searchUsers(searchQuery.value, 20);
      searchResults.value = response;
    } catch (err) {
      console.error('Search users error:', err);
    } finally {
      isLoading.value = false;
    }
  };

const startChat = (userId: number, userName: string) => {
  const existingConv = conversations.value.find(c => c.id === userId);
  if (existingConv) {
    selectConversation(userId);
  } else {
    conversations.value.unshift({
      id: userId,
      name: userName,
      avatar: '',
      lastMessage: '',
      time: '刚刚',
      unread: 0
    });
    selectConversation(userId);
  }
  showSearch.value = false;
  searchQuery.value = '';
};

const clearSearch = () => {
  searchQuery.value = '';
  searchResults.value = [];
};

onMounted(() => {
  loadConversations().then(() => {
    const userId = route.query.user_id;
    if (userId) {
      const userIdNum = parseInt(userId as string);
      if (!isNaN(userIdNum) && userIdNum > 0) {
        openChatWithUser(userIdNum);
      }
    }
  });
  initSSE();
});

onUnmounted(() => {
  closeSSE();
});
</script>

<template>
  <div class="message-page">
    <div v-if="!currentConversation" class="conversation-list">
      <header class="message-header">
        <button class="back-btn" @click="router.back()">
          <ArrowLeft :size="24" />
        </button>
        <h1 class="message-title">消息</h1>
        <button class="search-btn" @click="showSearch = !showSearch">
          <Search :size="24" />
        </button>
      </header>

      <div v-if="showSearch" class="search-box-container">
        <div class="search-box">
          <Search :size="20" />
          <input 
            v-model="searchQuery" 
            type="text" 
            placeholder="搜索用户..." 
            class="search-input"
            @input="searchUsers"
          />
          <button v-if="searchQuery" class="clear-btn" @click="clearSearch">
            <X :size="18" />
          </button>
        </div>
        
        <div v-if="searchResults.length > 0" class="search-results">
          <div 
            v-for="user in searchResults" 
            :key="user.id"
            class="search-item"
            @click="startChat(user.id, user.username)"
          >
            <div class="user-avatar">
              <User :size="24" />
            </div>
            <span class="user-name">{{ user.username }}</span>
          </div>
        </div>
      </div>

      <div class="conversations-container">
        <div 
          v-for="conv in conversations" 
          :key="conv.id"
          class="conversation-item"
          @click="selectConversation(conv.id)"
        >
          <div class="conv-avatar">
            <User :size="24" />
          </div>
          <div class="conv-info">
            <div class="conv-header">
              <span class="conv-name">{{ conv.name }}</span>
              <span class="conv-time">{{ conv.time }}</span>
            </div>
            <p class="conv-last-message">{{ conv.last_message || '开始聊天' }}</p>
          </div>
          <div v-if="conv.unread > 0" class="unread-badge">{{ conv.unread }}</div>
        </div>

        <div v-if="conversations.length === 0 && !showSearch" class="empty-state">
          <User :size="48" />
          <p>暂无消息</p>
          <button class="start-chat-btn" @click="showSearch = true">
            搜索用户开始聊天
          </button>
        </div>
      </div>
    </div>

    <div v-else class="chat-view">
      <header class="chat-header">
        <button class="back-btn" @click="currentConversation = null">
          <ArrowLeft :size="24" />
        </button>
        <div class="chat-title">
          <span class="chat-name">{{ conversations.find(c => c.id === currentConversation)?.name }}</span>
        </div>
        <button class="more-btn">
          <MoreVertical :size="24" />
        </button>
      </header>

      <div class="messages-container">
        <div 
          v-for="msg in messages" 
          :key="msg.id"
          class="message-item"
          :class="{ 'is-me': msg.sender_id === userStore.accountId }"
        >
          <div class="message-content">
            <p>{{ msg.content }}</p>
            <span class="message-time">{{ msg.created_at || msg.time }}</span>
          </div>
        </div>

        <div v-if="messages.length === 0" class="empty-chat">
          <p>开始与对方聊天吧</p>
        </div>
      </div>

      <div class="input-container">
        <input 
          v-model="newMessage" 
          type="text" 
          placeholder="输入消息..." 
          class="message-input"
          @keyup.enter="sendMessage"
        />
        <button class="send-btn" @click="sendMessage">
          <Send :size="20" />
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.message-page {
  min-height: 100vh;
  background: #000;
}

.conversation-list {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.message-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 56px;
  background: rgba(0, 0, 0, 0.9);
  display: flex;
  align-items: center;
  padding: 0 16px;
  z-index: 100;
}

.back-btn, .search-btn {
  background: none;
  border: none;
  color: #fff;
  cursor: pointer;
  padding: 8px;
}

.message-title {
  flex: 1;
  text-align: center;
  margin: 0;
  font-size: 18px;
}

.search-box-container {
  position: fixed;
  top: 56px;
  left: 0;
  right: 0;
  background: #000;
  padding: 12px 16px;
  z-index: 99;
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

.search-results {
  margin-top: 12px;
  max-height: 300px;
  overflow-y: auto;
}

.search-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  cursor: pointer;
  border-radius: 8px;
}

.search-item:hover {
  background: rgba(255, 255, 255, 0.05);
}

.user-avatar {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #ff0050, #ff2d55);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.user-name {
  font-size: 16px;
}

.conversations-container {
  flex: 1;
  padding: 72px 16px 16px;
}

.conversation-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  margin-bottom: 12px;
  cursor: pointer;
}

.conv-avatar {
  width: 50px;
  height: 50px;
  background: linear-gradient(135deg, #ff0050, #ff2d55);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.conv-info {
  flex: 1;
}

.conv-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 4px;
}

.conv-name {
  font-weight: bold;
  font-size: 16px;
}

.conv-time {
  color: #999;
  font-size: 12px;
}

.conv-last-message {
  margin: 0;
  color: #999;
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.unread-badge {
  background: #ff2d55;
  color: #fff;
  font-size: 12px;
  font-weight: bold;
  padding: 2px 8px;
  border-radius: 10px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
}

.empty-state svg {
  color: #333;
  margin-bottom: 16px;
}

.empty-state p {
  color: #999;
  margin: 0 0 16px 0;
}

.start-chat-btn {
  background: #ff2d55;
  border: none;
  border-radius: 20px;
  color: #fff;
  padding: 10px 24px;
  font-size: 14px;
  cursor: pointer;
}

.chat-view {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.chat-header {
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

.chat-title {
  flex: 1;
  text-align: center;
}

.chat-name {
  font-size: 18px;
  font-weight: bold;
}

.more-btn {
  background: none;
  border: none;
  color: #fff;
  cursor: pointer;
  padding: 8px;
}

.messages-container {
  flex: 1;
  padding: 72px 16px 80px;
  overflow-y: auto;
}

.message-item {
  display: flex;
  margin-bottom: 16px;
}

.message-item.is-me {
  justify-content: flex-end;
}

.message-content {
  max-width: 70%;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 16px;
}

.is-me .message-content {
  background: #ff2d55;
}

.message-content p {
  margin: 0 0 4px 0;
  font-size: 14px;
}

.message-time {
  font-size: 10px;
  color: #999;
}

.is-me .message-time {
  color: rgba(255, 255, 255, 0.7);
}

.empty-chat {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.empty-chat p {
  color: #666;
}

.input-container {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  gap: 12px;
  padding: 12px 16px;
  background: rgba(0, 0, 0, 0.9);
}

.message-input {
  flex: 1;
  height: 40px;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 20px;
  padding: 0 16px;
  color: #fff;
  font-size: 14px;
  outline: none;
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
</style>