
document.getElementById('login-form').addEventListener('submit', function(event) {
    event.preventDefault();
    
    const username = document.getElementById('username').value.trim();
    const password = document.getElementById('password').value.trim();

    if (username && password) {
        alert(`Merci ${username}, vous etes connecté avec succès!`);
        
        document.getElementById('login-form').reset();
    } else {
        alert("Veuillez remplir tous les champs du formulaire.");
    }
});