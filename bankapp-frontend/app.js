const API_URL = "http://localhost:8080";

const loginSection = document.getElementById("login-section");
const registerSection = document.getElementById("register-section");
const accountsSection = document.getElementById("accounts-section");

const loginButton = document.getElementById("login-button");
const registerButton = document.getElementById("register-button");
const addAccountButton = document.getElementById("add-account-button");
const transferButton = document.getElementById("transfer-button");
const logoutButton = document.getElementById("logout-button");

const toRegisterLink = document.getElementById("to-register");
const toLoginLink = document.getElementById("to-login");

const loginError = document.getElementById("login-error");
const registerError = document.getElementById("register-error");

const accountsList = document.getElementById("accounts-list");

// Переключение между формами
toRegisterLink.addEventListener("click", () => {
  loginSection.classList.add("hidden");
  registerSection.classList.remove("hidden");
});

toLoginLink.addEventListener("click", () => {
  registerSection.classList.add("hidden");
  loginSection.classList.remove("hidden");
});

// Вход
loginButton.addEventListener("click", async () => {
  const username = document.getElementById("login-username").value;
  const password = document.getElementById("login-password").value;

  try {
    const response = await fetch(`${API_URL}/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });

    if (!response.ok) {
      const errorData = await response.json();
      loginError.textContent = errorData.error || "Ошибка входа";
      return;
    }

    const data = await response.json();
    localStorage.setItem("token", data.token);
    showAccountsSection();
  } catch (error) {
    loginError.textContent = "Сетевая ошибка";
  }
});

// Регистрация
registerButton.addEventListener("click", async () => {
  const username = document.getElementById("register-username").value;
  const password = document.getElementById("register-password").value;

  try {
    const response = await fetch(`${API_URL}/auth/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });

    if (!response.ok) {
      const errorData = await response.json();
      registerError.textContent = errorData.error || "Ошибка регистрации";
      return;
    }

    alert("Регистрация успешна!");
    registerSection.classList.add("hidden");
    loginSection.classList.remove("hidden");
  } catch (error) {
    registerError.textContent = "Сетевая ошибка";
  }
});

// Показ счетов
async function showAccountsSection() {
  loginSection.classList.add("hidden");
  accountsSection.classList.remove("hidden");

  const token = localStorage.getItem("token");

  try {
    const response = await fetch(`${API_URL}/accounts`, {
      headers: { Authorization: `Bearer ${token}` },
    });

    if (!response.ok) {
      throw new Error("Не удалось загрузить счета");
    }

    const accounts = await response.json();
    accountsList.innerHTML = "";
    accounts.forEach((account) => {
      const li = document.createElement("li");
      li.textContent = `ID: ${account.ID}, Баланс: ${account.Balance}`;
      accountsList.appendChild(li);
    });
  } catch (error) {
    console.error("Ошибка при загрузке счетов:", error);
  }
}

// Добавление счета
addAccountButton.addEventListener("click", async () => {
  const balance = parseFloat(document.getElementById("new-account-balance").value); // Преобразование строки в число
  const token = localStorage.getItem("token");

  // Проверка на корректность баланса
  if (isNaN(balance) || balance <= 0) {
    alert("Введите корректное числовое значение для баланса!");
    return;
  }

  try {
    const response = await fetch(`${API_URL}/accounts`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`, // Передача токена
      },
      body: JSON.stringify({ initial_balance: balance }), // Отправляем данные
    });

    if (!response.ok) {
      const errorData = await response.json();
      console.error("Ошибка добавления счёта:", errorData);
      alert(errorData.error || "Ошибка при добавлении счёта");
      return;
    }

    alert("Счёт успешно добавлен!");
    await fetchAndDisplayAccounts(); // Обновляем список счетов
  } catch (error) {
    console.error("Ошибка при добавлении счёта:", error);
    alert("Сетевая ошибка при добавлении счёта");
  }
});


async function fetchAndDisplayAccounts() {
  const token = localStorage.getItem("token");

  try {
    const response = await fetch(`${API_URL}/accounts`, {
      headers: { Authorization: `Bearer ${token}` },
    });

    if (!response.ok) {
      throw new Error("Не удалось загрузить счета");
    }

    const accounts = await response.json();

    // Очистка старого списка
    accountsList.innerHTML = "";

    // Отображение новых данных
    accounts.forEach((account) => {
      const li = document.createElement("li");
      li.textContent = `ID: ${account.ID}, Баланс: ${account.Balance}`;
      accountsList.appendChild(li);
    });
  } catch (error) {
    console.error("Ошибка при загрузке счетов:", error);
  }
}


// Переводы
transferButton.addEventListener("click", async () => {
  const fromAccount = document.getElementById("from-account-id").value;
  const toAccount = document.getElementById("to-account-id").value;
  const amount = document.getElementById("transfer-amount").value;
  const token = localStorage.getItem("token");

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
        category: "example",
      }),
    });

    if (!response.ok) {
      alert("Ошибка при переводе");
      return;
    }

    alert("Перевод выполнен!");
    showAccountsSection();
  } catch (error) {
    console.error("Ошибка при переводе:", error);
  }
});

// Выход
logoutButton.addEventListener("click", () => {
  localStorage.removeItem("token");
  accountsSection.classList.add("hidden");
  loginSection.classList.remove("hidden");
});
