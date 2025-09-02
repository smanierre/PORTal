<script setup lang="ts">
import { useAppStore } from '@/stores/app_store.ts';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

const appStore = useAppStore();
const router = useRouter();

onMounted(async () => {
  if (appStore.member) {
    await router.replace('/');
  }
});

const username = ref('');
const password = ref('');
const errorMessage = ref('');
const requestInProgress = ref(false);

async function login() {
  errorMessage.value = '';
  const id = setTimeout(() => {
    requestInProgress.value = true;
  }, 200);
  const res = await appStore.login(username.value, password.value);
  // Prevent it from setting the spinner if the request errors out quickly
  clearTimeout(id);
  if (res) {
    errorMessage.value = res.message;
    requestInProgress.value = false;
    return;
  }
  await router.push('/');
}
</script>

<template>
  <div class="login-container">
    <div class="login-form-container">
      <p class="login-form-header">{{ appStore.organization }} PORTal Login</p>
      <form @submit.prevent="login">
        <label>
          <input
            v-model="username"
            class="login-input"
            name="username"
            placeholder="Username"
            style="margin-top: 2.5rem"
            type="text"
          />
        </label>
        <label>
          <input
            v-model="password"
            class="login-input"
            name="password"
            placeholder="Password"
            type="password"
          />
        </label>
        <div class="login-form-submit-container">
          <button
            :class="{ 'form-submit': true, 'submit-active': requestInProgress }"
            class="form-submit"
            :disabled="username && password ? false : true"
          >
            <span class="form-submit--text">Login</span>
            <span class="form-submit--spinner"
              ><svg
                class="spinner"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  d="M10.14,1.16a11,11,0,0,0-9,8.92A1.59,1.59,0,0,0,2.46,12,1.52,1.52,0,0,0,4.11,10.7a8,8,0,0,1,6.66-6.61A1.42,1.42,0,0,0,12,2.69h0A1.57,1.57,0,0,0,10.14,1.16Z"
                  class="spinner_P7sC"
                /></svg
            ></span>
          </button>
        </div>
      </form>
      <Transition>
        <div v-if="errorMessage" class="login-error">
          {{ errorMessage }}
        </div>
      </Transition>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  width: 100vw;
  height: 100vh;
  justify-content: center;
  align-items: center;
}

.login-form-container {
  width: 24rem;
  height: fit-content;
  padding-bottom: 2rem;
  background-color: hsl(var(--background));
  position: relative;
}

.login-form-header {
  text-align: center;
  font-size: 1.25rem;
  line-height: 1.75rem;
  font-weight: 700;
  margin-top: 1rem;
  color: hsl(var(--primary));
}

.login-input {
  width: 80%;
  height: 2.5rem;
  display: block;
  margin-left: auto;
  margin-right: auto;
  margin-top: 1.25rem;
}

.login-form-submit-container {
  display: flex;
  justify-content: center;
  margin-top: 1rem;
}

.v-enter-active,
.v-leave-active {
  transition: opacity 0.5s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
  translate: 0 50px;
}

.login-error {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  color: rgb(var(--red));
  border: 1px solid rgb(var(--red-dark));
  text-align: center;
  line-height: 1.25rem;
  font-size: 1rem;
  width: 80%;
  margin-top: 1.25rem;
  padding: 0.5rem 1rem;
  transition-property: all;
  transition-duration: 500ms;
  border-radius: 0.25rem;
}

.spinner_P7sC {
  transform-origin: center;
  animation: spinner_svv2 0.75s infinite linear;
}

@keyframes spinner_svv2 {
  100% {
    transform: rotate(360deg);
  }
}

.form-submit {
  display: grid;
  grid-template-areas: 'stack';

  .form-submit--spinner {
    opacity: 0;
    grid-area: stack;
  }

  .form-submit--text {
    grid-area: stack;
  }
}

.submit-active {
  .form-submit--spinner {
    opacity: 1;
  }

  .form-submit--text {
    opacity: 0;
  }
}
</style>
