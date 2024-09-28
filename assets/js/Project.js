
function validateProjectForm() {
    const title = document.getElementById('project-title').value.trim();
    const description = document.getElementById('project-description').value.trim();
    const image = document.getElementById('project-image').value.trim();
    const link = document.getElementById('project-link').value.trim();

    if (!title || !description || !image || !link) {
        alert("Tous les champs du formulaire sont obligatoires.");
        return false;
    }
    return true;
}


function showConfirmationMessage() {
    const confirmation = document.getElementById('confirmation-message');
    if (confirmation) {
        confirmation.style.display = 'block';
        setTimeout(() => {
            confirmation.style.display = 'none';
        }, 3000);
    }
}

window.onload = function () {
    showConfirmationMessage();
};

document.addEventListener('DOMContentLoaded', function () {
    const addProjectForm = document.getElementById('add-project-form');
    if (addProjectForm) {
        addProjectForm.onsubmit = validateProjectForm;
    }
});


