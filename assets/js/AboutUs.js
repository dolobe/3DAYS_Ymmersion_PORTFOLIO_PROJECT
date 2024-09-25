// assets/js/AboutUs.js
// Ajout d'un effet de surbrillance sur les cartes de compétences
document.querySelectorAll('.skill-card').forEach(card => {
    card.addEventListener('mouseover', () => {
        card.style.backgroundColor = '#d4edda';
        card.style.transition = 'background-color 0.3s ease';
    });

    card.addEventListener('mouseout', () => {
        card.style.backgroundColor = '#e9ecef';
    });
});
