const submitGuessBtn = document.getElementById('submit-guess');
const result = document.getElementById('result');

let isLoggedIn = false;
if (localStorage.getItem('token')) {
	isLoggedIn = true;
}
if (isLoggedIn) {
	document.getElementById('guess-container').style.display = 'block';
} else {
	window.location.href = '../static/login.html';
}

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

