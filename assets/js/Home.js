// Smooth scroll to sections
function scrollToSection(sectionId) {
    document.getElementById(sectionId).scrollIntoView({ behavior: 'smooth' });
}

// Handling the contact form submission
document.getElementById('contact-form').addEventListener('submit', function(event) {
    event.preventDefault();
    
    const name = document.getElementById('name').value;
    const email = document.getElementById('email').value;
    const message = document.getElementById('message').value;

    if (name && email && message) {
        alert(`Merci, ${name}! Votre message a bien été envoyé.`);
        // Clear form fields
        document.getElementById('contact-form').reset();
    } else {
        alert("Veuillez remplir tous les champs du formulaire.");
    }
});
