// app.js

// Элементы DOM
const loginForm = document.getElementById('login-form');
const accountsSection = document.getElementById('accounts-section');
const loginButton = document.getElementById('login-button');
const logoutButton = document.getElementById('logout-button');
const usernameInput = document.getElementById('username');
const passwordInput = document.getElementById('password');
const loginError = document.getElementById('login-error');
const accountsList = document.getElementById('accounts-list');

// URL вашего бэкенд API
const API_URL = 'http://localhost:8080';

// Проверка наличия токена при загрузке страницы
window.onload = () => {
  const token = localStorage.getItem('token');
  if (token) {
    showAccounts();
  }
};

// Обработчик входа
loginButton.addEventListener('click', async () => {
  const username = usernameInput.value;
  const password = passwordInput.value;

  try {
    const response = await fetch(`${API_URL}/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ username, password })
    });

    if (!response.ok) {
      const errorData = await response.json();
      loginError.textContent = errorData.error || 'Ошибка авторизации';
      return;
    }

    const data = await response.json();
    localStorage.setItem('token', data.token);
    showAccounts();
  } catch (error) {
    loginError.textContent = 'Сетевая ошибка';
    console.error('Ошибка при авторизации:', error);
  }
});

// Обработчик выхода
logoutButton.addEventListener('click', () => {
  localStorage.removeItem('token');
  loginForm.style.display = 'block';
  accountsSection.style.display = 'none';
});

// Функция для отображения счетов
async function showAccounts() {
  const token = localStorage.getItem('token');

  try {
    const response = await fetch(`${API_URL}/accounts`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });

    if (!response.ok) {
      localStorage.removeItem('token');
      loginForm.style.display = 'block';
      accountsSection.style.display = 'none';
      return;
    }

    const data = await response.json();
    accountsList.innerHTML = '';
    data.accounts.forEach(account => {
      const li = document.createElement('li');
      li.textContent = `ID: ${account.ID}, Баланс: ${account.Balance}₽`;
      accountsList.appendChild(li);
    });

    loginForm.style.display = 'none';
    accountsSection.style.display = 'block';
  } catch (error) {
    console.error('Ошибка при загрузке счетов:', error);
  }
}
