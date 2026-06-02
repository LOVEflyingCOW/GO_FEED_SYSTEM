import axios from 'axios'
import type { LoginResponse, FeedResponse, Comment, LikeResponse, Profile, Message, Video } from '@/types'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('account_id')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export const accountAPI = {
  register(username: string, password: string) {
    return api.post('/account/register', { username, password })
  },
  login(username: string, password: string): Promise<LoginResponse> {
    return api.post('/account/login', { username, password })
  },
  refresh(refreshToken: string) {
    return api.post('/account/refresh', { refresh_token: refreshToken })
  },
  logout() {
    return api.post('/account/logout')
  },
  getAccount(id: number): Promise<{ id: number; username: string; avatar_url?: string }> {
    return api.get(`/account/${id}`)
  },
  searchUsers(keyword: string, limit?: number): Promise<{ id: number; username: string; avatar_url?: string }[]> {
    return api.get('/account/search', { params: { keyword, limit } })
  },
  rename(newUsername: string) {
    return api.put('/account/rename', { new_username: newUsername })
  },
  updateProfile(username?: string, bio?: string) {
    return api.put('/account/profile', { username, bio })
  },
  uploadAvatar(formData: FormData) {
    return api.post('/account/avatar', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  }
}

export const feedAPI = {
  getLatest(cursor?: number, limit?: number): Promise<FeedResponse> {
    return api.get('/feed', { params: { cursor, limit, type: 'latest' } })
  },
  getHot(limit?: number): Promise<FeedResponse> {
    return api.get('/feed/hot', { params: { limit } })
  },
  getFollowing(cursor?: number, limit?: number): Promise<FeedResponse> {
    return api.get('/feed/following', { params: { cursor, limit } })
  },
  getTag(tag: string, cursor?: number, limit?: number): Promise<FeedResponse> {
    return api.get(`/feed/tag/${tag}`, { params: { cursor, limit } })
  },
  search(keyword: string, limit?: number): Promise<FeedResponse> {
    return api.get('/feed/search', { params: { keyword, limit } })
  }
}

export const videoAPI = {
  getVideo(videoId: number) {
    return api.get(`/video/${videoId}`)
  },
  getUserVideos(accountId: number, page?: number, limit?: number): Promise<{ videos: Video[]; total: number }> {
    return api.get(`/video/user/${accountId}`, { params: { page, limit } })
  },
  upload(formData: FormData) {
    return api.post('/video/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },
  delete(videoId: number) {
    return api.delete(`/video/${videoId}`)
  },
  reportView(videoId: number) {
    return api.post(`/video/${videoId}/view`)
  }
}

export const likeAPI = {
  like(videoId: number): Promise<LikeResponse> {
    return api.post(`/like/${videoId}`)
  },
  unlike(videoId: number): Promise<LikeResponse> {
    return api.delete(`/like/${videoId}`)
  },
  getStatus(videoId: number): Promise<LikeResponse> {
    return api.get(`/like/${videoId}`)
  },
  getLikedVideos(cursor?: number, limit?: number): Promise<{ videos: Video[]; total: number }> {
    return api.get('/like/list', { params: { cursor, limit } })
  }
}

export const commentAPI = {
  create(videoId: number, content: string, replyTo?: number) {
    return api.post(`/comment/${videoId}`, { content, reply_to: replyTo })
  },
  list(videoId: number, page?: number, limit?: number): Promise<{ comments: Comment[]; total: number }> {
    return api.get(`/comment/${videoId}/list`, { params: { page, limit } })
  },
  delete(commentId: number) {
    return api.delete(`/comment/${commentId}`)
  }
}

export const socialAPI = {
  follow(targetId: number) {
    return api.post(`/social/follow/${targetId}`)
  },
  unfollow(targetId: number) {
    return api.delete(`/social/unfollow/${targetId}`)
  },
  getProfile(accountId: number): Promise<Profile> {
    return api.get(`/social/profile/${accountId}`)
  },
  getFollowers(targetId: number, page?: number, limit?: number) {
    return api.get(`/social/followers/${targetId}`, { params: { page, limit } })
  },
  getFollowing(accountId: number, page?: number, limit?: number) {
    return api.get(`/social/following/${accountId}`, { params: { page, limit } })
  }
}

export const messageAPI = {
  send(receiverId: number, content: string) {
    return api.post('/message/send', { receiver_id: receiverId, content })
  },
  getConversations(): Promise<{ user_id: number; username: string; avatar_url?: string; last_message: string; unread_count: number; updated_at: string }[]> {
    return api.get('/message/conversations')
  },
  getMessages(otherId: number, page?: number, limit?: number): Promise<{ messages: Message[]; total: number }> {
    return api.get(`/message/${otherId}`, { params: { page, limit } })
  }
}