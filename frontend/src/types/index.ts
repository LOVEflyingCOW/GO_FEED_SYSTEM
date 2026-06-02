export interface Account {
  id: number
  username: string
  avatar_url?: string
  bio?: string
}

export interface LoginResponse {
  token: string
  refresh_token: string
  account_id: number
  username: string
}

export interface Video {
  id: number
  account_id: number
  username: string
  title: string
  video_path: string
  cover_path: string
  play_url: string
  cover_url: string
  duration: number
  description?: string
  tags: string
  view_count: number
  like_count: number
  comment_count: number
  created_at: string
}

export interface FeedItem {
  video_id: number
  account_id: number
  username: string
  avatar_url?: string
  title: string
  video_url: string
  cover_url: string
  duration: number
  description?: string
  tags: string
  view_count: number
  like_count: number
  comment_count: number
  created_at: string
  is_liked?: boolean
  is_following?: boolean
}

export interface FeedResponse {
  items: FeedItem[]
  has_more: boolean
  next: number
}

export interface Comment {
  id: number
  account_id: number
  username: string
  avatar_url?: string
  video_id: number
  content: string
  reply_to?: number
  created_at: string
}

export interface LikeResponse {
  video_id: number
  like_count: number
  is_liked: boolean
}

export interface Profile {
  account_id: number
  username: string
  avatar_url?: string
  bio?: string
  follower_count: number
  following_count: number
  is_followed?: boolean
}

export interface Message {
  id: number
  sender_id: number
  sender_name: string
  content: string
  is_read: boolean
  created_at: string
}

export interface Notification {
  type: string
  from_id: number
  from: string
  content: unknown
}