// assets/js/Contact.js
document.getElementById('contact-form').addEventListener('submit', function(event) {
    event.preventDefault();
    
    const name = document.getElementById('name').value.trim();
    const email = document.getElementById('email').value.trim();
    const subject = document.getElementById('subject').value.trim();
    const message = document.getElementById('message').value.trim();

    if (name && email && subject && message) {
        // Afficher un message de confirmation
        alert(`Merci ${name}, votre message a été envoyé avec succès!`);
        
        // Effacer les champs du formulaire après l'envoi
        document.getElementById('contact-form').reset();
    } else {
        alert("Veuillez remplir tous les champs du formulaire.");
    }
});
