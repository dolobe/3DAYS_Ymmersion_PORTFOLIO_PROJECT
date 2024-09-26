document.getElementById('login-form').addEventListener('submit', function(event) {
    event.preventDefault();

    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    if (username === 'admin' && password === 'password') {
        alert('Connexion réussie !');
        
        window.location.href = 'admin.html';
    } else {
        document.getElementById('message').textContent = 'Nom d\'utilisateur ou mot de passe incorrect.';
    }
});
