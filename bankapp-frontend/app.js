const API_URL = "http://localhost:8080";

// Sections
const loginSection = document.getElementById("login-section");
const registerSection = document.getElementById("register-section");
const accountsSection = document.getElementById("accounts-section");

// Buttons
const loginButton = document.getElementById("login-button");
const registerButton = document.getElementById("register-button");
const addAccountButton = document.getElementById("add-account-button");
const transferButton = document.getElementById("transfer-button");
const logoutButton = document.getElementById("logout-button");

// Navigation Links
const toRegisterLink = document.getElementById("to-register");
const toLoginLink = document.getElementById("to-login");

// Error messages
const loginError = document.getElementById("login-error");
const registerError = document.getElementById("register-error");

// Account list
const accountsList = document.getElementById("accounts-list");

// Helper Functions
function switchSection(from, to) {
  from.classList.remove("visible");
  from.classList.add("hidden");
  setTimeout(() => {
    from.style.display = "none";
    to.style.display = "flex";
    to.classList.remove("hidden");
    to.classList.add("visible");
  }, 500); // Время соответствует длительности transition в CSS
}

async function fetchWithAuth(endpoint, options = {}) {
  const token = localStorage.getItem("token");
  const headers = { Authorization: `Bearer ${token}`, ...options.headers };
  const response = await fetch(`${API_URL}${endpoint}`, { ...options, headers });
  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.error || "Произошла ошибка");
  }
  return response.json();
}

// Navigation between forms
toRegisterLink.addEventListener("click", () => switchSection(loginSection, registerSection));
toLoginLink.addEventListener("click", () => switchSection(registerSection, loginSection));

// Login
loginButton.addEventListener("click", async () => {
  const username = document.getElementById("login-username").value;
  const password = document.getElementById("login-password").value;

  try {
    const data = await fetch(`${API_URL}/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    }).then((res) => res.json());

    if (!data.token) throw new Error("Ошибка авторизации");

    localStorage.setItem("token", data.token);
    switchSection(loginSection, accountsSection);
    await fetchAndDisplayAccounts();
  } catch (error) {
    loginError.textContent = error.message;
  }
});

// Registration
registerButton.addEventListener("click", async () => {
  const username = document.getElementById("register-username").value;
  const password = document.getElementById("register-password").value;

  try {
    await fetch(`${API_URL}/auth/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });
    alert("Регистрация успешна!");
    switchSection(registerSection, loginSection);
  } catch (error) {
    registerError.textContent = error.message;
  }
});

// Fetch and display accounts
async function fetchAndDisplayAccounts() {
  const token = localStorage.getItem("token");

  try {
    const response = await fetch(`${API_URL}/accounts`, {
      headers: { Authorization: `Bearer ${token}` },
    });

    if (!response.ok) {
      throw new Error(`Ошибка: ${response.status}`);
    }

    const responseData = await response.json();

    // Логирование данных для проверки
    console.log("Полученные данные:", responseData);

    const accounts = responseData.accounts;

    // Проверяем, что `accounts` — массив
    if (!Array.isArray(accounts)) {
      console.error("Ошибка формата: `accounts` не является массивом", accounts);
      throw new Error("Неверный формат данных, ожидается массив счетов");
    }

    // Очистка текущего списка и отображение новых данных
    accountsList.innerHTML = "";
    accounts.forEach((account) => {
      const li = document.createElement("li");
      li.textContent = `ID: ${account.ID}, Баланс: ${account.Balance}`;
      accountsList.appendChild(li);
    });

    console.log("Счета успешно отображены");
  } catch (error) {
    console.error("Ошибка при загрузке счетов:", error);
    alert("Не удалось загрузить счета. Проверьте подключение или повторите позже.");
  }
}

async function showAccountsSection() {
  // Скрываем секцию входа и показываем секцию счетов
  loginSection.classList.add("hidden");
  accountsSection.classList.remove("hidden");

  const token = localStorage.getItem("token"); // Получаем токен из локального хранилища

  try {
    // Запрашиваем данные счетов с сервера
    const response = await fetch(`${API_URL}/accounts`, {
      headers: { Authorization: `Bearer ${token}` }, // Передаем токен в заголовках
    });

    if (!response.ok) {
      throw new Error("Не удалось загрузить счета");
    }

    const responseData = await response.json(); // Парсим ответ сервера

    // Проверяем, что ответ содержит массив счетов
    const accounts = responseData.accounts || []; // Предотвращаем ошибку, если accounts отсутствует
    if (!Array.isArray(accounts)) {
      throw new Error("Неверный формат данных счетов");
    }

    // Очищаем текущий список счетов на странице
    accountsList.innerHTML = "";

    // Перебираем счета и отображаем их
    accounts.forEach((account) => {
      const li = document.createElement("li");
      li.textContent = `ID: ${account.ID}, Баланс: ${account.Balance}`;
      accountsList.appendChild(li);
    });
  } catch (error) {
    // Логируем ошибки в консоль для отладки
    console.error("Ошибка при загрузке счетов:", error);
    alert("Не удалось загрузить счета. Проверьте подключение или повторите позже.");
  }
}


// Add account
addAccountButton.addEventListener("click", async () => {
  const balance = parseFloat(document.getElementById("new-account-balance").value);

  if (isNaN(balance) || balance <= 0) {
    alert("Введите корректное числовое значение для баланса!");
    return;
  }

  try {
    await fetchWithAuth("/accounts", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ initial_balance: balance }),
    });
    alert("Счёт успешно добавлен!");
    await fetchAndDisplayAccounts();
  } catch (error) {
    console.error("Ошибка при добавлении счёта:", error);
    alert(error.message);
  }
});

// Transfer funds
// Переводы
transferButton.addEventListener("click", async () => {
  const fromAccount = parseInt(document.getElementById("from-account-id").value, 10);
  const toAccount = parseInt(document.getElementById("to-account-id").value, 10);
  const amount = parseFloat(document.getElementById("transfer-amount").value);
  const token = localStorage.getItem("token");

  // Проверка на валидность данных
  if (isNaN(fromAccount) || isNaN(toAccount) || isNaN(amount) || amount <= 0) {
    alert("Введите корректные значения для перевода.");
    return;
  }

  try {
    const response = await fetch(`${API_URL}/transfer`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        from_account_id: fromAccount,
        to_account_id: toAccount,
        amount: amount,
        category: "example", // Пример категории, если требуется
      }),
    });

    if (!response.ok) {
      const errorData = await response.json();
      console.error("Ошибка при переводе:", errorData);
      alert(errorData.error || "Ошибка при переводе.");
      return;
    }

    alert("Перевод выполнен успешно!");
    await fetchAndDisplayAccounts(); // Обновляем список счетов
  } catch (error) {
    console.error("Ошибка при переводе:", error);
    alert("Сетевая ошибка при переводе.");
  }
});


// Logout
logoutButton.addEventListener("click", () => {
  localStorage.removeItem("token");
  switchSection(accountsSection, loginSection);
});