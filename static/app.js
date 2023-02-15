const loginForm = document.getElementById('login-form');
const loginContainer = document.getElementById('login-container');
const guessContainer = document.getElementById('guess-container');
const submitGuessBtn = document.getElementById('submit-guess');
const result = document.getElementById('result');

loginForm.addEventListener('submit', (event) => {
	event.preventDefault();

	const username = document.getElementById('username').value;
	const password = document.getElementById('password').value;

	fetch('/login', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({username, password})
	})
	.then(response => response.json())
	.then(data => {
		if (data.token) {
			localStorage.setItem('token', data.token);
			loginContainer.style.display = 'none';
			guessContainer.style.display = 'block';
		}
	})
	.catch(error => console.error(error));
});

submitGuessBtn.addEventListener('click', () => {
	const guess = document.getElementById('guess').value;
	const token = localStorage.getItem('token');

	fetch('/guess', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
			'Authorization': `Bearer ${token}`
		},
		body: JSON.stringify({guess})
	})
	.then(response => response.json())
	.then(data => {
		if (data.result) {
			result.innerHTML = `Your guess was ${data.result}.`;
			result.style.display = 'block';
		}
	})
	.catch(error => console.error(error));
});
