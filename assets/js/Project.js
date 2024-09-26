document.addEventListener("DOMContentLoaded", () => {
    const form = document.getElementById("project-form");

    async function fetchProjects() {
        console.log("Tentative de récupération des projets...");
        try {
            const response = await fetch("/api/projects");
            console.log("Réponse reçue de l'API :", response);
            if (!response.ok) {
                throw new Error("Erreur lors de la récupération des projets");
            }
            const projects = await response.json();
            console.log("Projets récupérés :", projects);
            projects.forEach(project => addProjectToList(project));
        } catch (error) {
            console.error("Erreur lors de la récupération des projets :", error);
        }
    }
    
    function addProjectToList(project) {
        console.log("Ajout du projet à la liste :", project);
        const projectList = document.getElementById("project-list");
        const projectItem = document.createElement("div");
        projectItem.className = "project-item card";

        projectItem.style.border = "2px solid grey";
        projectItem.style.borderRadius = "10px";
        projectItem.style.width = "300px";
        projectItem.style.height = "400px";
        projectItem.style.display = "flex";
        projectItem.style.flexDirection = "column";
        projectItem.style.justifyContent = "space-between";
        projectItem.style.alignItems = "center";

        projectItem.innerHTML = `
            <img src="${project.image}" alt="${project.title}" class="project-image" style="width: 100%; height: auto; border-radius: 5px;">
            <div class="project-details">
                <h3>${project.title}</h3>
                <p>${project.description}</p>
                <a href="${project.link}" target="_blank" class="project-link">Voir le Projet</a>
            </div>
        `;

        const projectLink = projectItem.querySelector(".project-link");

        projectLink.style.backgroundColor = "#4CAF50";
        projectLink.style.color = "white";
        projectLink.style.padding = "10px 20px";
        projectLink.style.textAlign = "center";
        projectLink.style.textDecoration = "none";
        projectLink.style.display = "inline-block";
        projectLink.style.borderRadius = "5px";
        projectLink.style.transition = "background-color 0.3s";
        projectLink.style.marginBottom = "10px";

        projectLink.addEventListener("mouseenter", () => {
            projectLink.style.backgroundColor = "#45a049";
        });
        projectLink.addEventListener("mouseleave", () => {
            projectLink.style.backgroundColor = "#4CAF50";
        });

        projectList.appendChild(projectItem);
        console.log("Projet ajouté au DOM avec succès.");
    }

    fetchProjects();

    form.addEventListener("submit", async (event) => {
        event.preventDefault();
        console.log("Formulaire soumis, tentative d'ajout du projet...");

        const projectTitle = document.getElementById("project-title").value;
        const projectDescription = document.getElementById("project-description").value;
        const projectImage = document.getElementById("project-image").value;
        const projectLink = document.getElementById("project-link").value;

        const project = {
            title: projectTitle,
            description: projectDescription,
            image: projectImage,
            link: projectLink
        };

        console.log("Données du projet à envoyer :", project);

        try {
            const response = await fetch("/api/projects", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(project)
            });

            console.log("Réponse reçue de l'ajout de projet :", response);
            if (!response.ok) {
                throw new Error("Erreur lors de l'ajout du projet");
            }

            const newProject = await response.json();
            console.log("Projet ajouté avec succès :", newProject);
            addProjectToList(newProject);

            form.reset();
            console.log("Formulaire réinitialisé avec succès.");
        } catch (error) {
            console.error("Erreur lors de l'ajout du projet :", error);
            alert("Une erreur s'est produite lors de l'ajout du projet.");
        }
    });
});
