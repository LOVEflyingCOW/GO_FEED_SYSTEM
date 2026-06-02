<script setup lang="ts">import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/stores/user';
import { User, Lock, Eye, EyeOff } from 'lucide-vue-next';
const router = useRouter();
const userStore = useUserStore();
const username = ref('');
const password = ref('');
const showPassword = ref(false);
const error = ref('');
const loading = ref(false);
const handleLogin = async () => {
 if (!username.value || !password.value) {
 error.value = '请输入用户名和密码';
 return;
 }
 loading.value = true;
 try {
 await userStore.login(username.value, password.value);
 await router.push('/');
 }
 catch (err) {
 error.value = '用户名或密码错误';
 }
 finally {
 loading.value = false;
 }
};
</script>

<template>
  <div class="login-page">
    <div class="login-container">
      <div class="logo-section">
        <div class="logo">
          <span class="logo-icon">🎵</span>
        </div>
        <h1 class="app-name">抖音</h1>
        <p class="app-slogan">记录美好生活</p>
      </div>

      <div class="form-section">
        <div class="input-group">
          <User class="input-icon" :size="20" />
          <input 
            v-model="username" 
            type="text" 
            placeholder="用户名" 
            class="input-field"
            @keyup.enter="handleLogin"
          />
        </div>

        <div class="input-group">
          <Lock class="input-icon" :size="20" />
          <input 
            v-model="password" 
            :type="showPassword ? 'text' : 'password'" 
            placeholder="密码" 
            class="input-field"
            @keyup.enter="handleLogin"
          />
          <button class="toggle-password" @click="showPassword = !showPassword">
            <Eye v-if="!showPassword" :size="20" />
            <EyeOff v-else :size="20" />
          </button>
        </div>

        <p v-if="error" class="error-message">{{ error }}</p>

        <button 
          class="login-btn" 
          :disabled="loading"
          @click="handleLogin"
        >
          <span v-if="loading">登录中...</span>
          <span v-else>登录</span>
        </button>

        <p class="register-link">
          还没有账号？
          <span class="link" @click="router.push('/register')">立即注册</span>
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(180deg, #1a1a2e 0%, #16213e 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.login-container {
  width: 100%;
  max-width: 400px;
  text-align: center;
}

.logo-section {
  margin-bottom: 40px;
}

.logo {
  width: 100px;
  height: 100px;
  margin: 0 auto;
  background: linear-gradient(135deg, #ff0050, #ff2d55);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
}

.logo-icon {
  font-size: 48px;
}

.app-name {
  font-size: 48px;
  font-weight: bold;
  background: linear-gradient(135deg, #ff0050, #ff2d55);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0 0 8px 0;
}

.app-slogan {
  color: #999;
  font-size: 14px;
  margin: 0;
}

.form-section {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 20px;
  padding: 30px;
  backdrop-filter: blur(10px);
}

.input-group {
  position: relative;
  margin-bottom: 20px;
}

.input-icon {
  position: absolute;
  left: 15px;
  top: 50%;
  transform: translateY(-50%);
  color: #666;
}

.input-field {
  width: 100%;
  height: 48px;
  padding: 0 15px 0 45px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 24px;
  color: #fff;
  font-size: 16px;
  outline: none;
  transition: border-color 0.3s;
}

.input-field:focus {
  border-color: #ff2d55;
}

.input-field::placeholder {
  color: #666;
}

.toggle-password {
  position: absolute;
  right: 15px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: #666;
  cursor: pointer;
}

.error-message {
  color: #ff4444;
  font-size: 14px;
  margin: -10px 0 15px 0;
}

.login-btn {
  width: 100%;
  height: 48px;
  background: linear-gradient(135deg, #ff0050, #ff2d55);
  border: none;
  border-radius: 24px;
  color: #fff;
  font-size: 16px;
  font-weight: bold;
  cursor: pointer;
  transition: opacity 0.3s;
  margin-bottom: 20px;
}

.login-btn:hover:not(:disabled) {
  opacity: 0.9;
}

.login-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.register-link {
  color: #999;
  font-size: 14px;
}

.link {
  color: #ff2d55;
  cursor: pointer;
}

.link:hover {
  text-decoration: underline;
}
</style>