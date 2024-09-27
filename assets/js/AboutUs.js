
document.querySelectorAll('.skill-card').forEach(card => {
    card.addEventListener('mouseover', () => {
        card.style.backgroundColor = '#d4edda';
        card.style.transition = 'background-color 0.3s ease';
    });

    card.addEventListener('mouseout', () => {
        card.style.backgroundColor = '#e9ecef';
    });
});

document.addEventListener("DOMContentLoaded", function () {
    const aboutForm = document.getElementById("aboutForm");
    const nameField = document.getElementById("name");
    const contentField = document.getElementById("content");

    // Récupérer les données actuelles
    fetch("/About")
        .then(response => response.json())
        .then(data => {
            nameField.value = data.name;
            contentField.value = data.content;
        });

    aboutForm.addEventListener("submit", function (event) {
        event.preventDefault();

        const aboutData = {
            name: nameField.value,
            content: contentField.value
        };

        fetch("/About", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(aboutData)
        })
            .then(response => {
                if (response.ok) {
                    alert("Informations mises à jour avec succès !");
                } else {
                    alert("Erreur lors de la mise à jour des informations.");
                }
            })
            .catch(error => console.error("Erreur :", error));
    });
});

