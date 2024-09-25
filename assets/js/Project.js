document.addEventListener('DOMContentLoaded', function () {
    const projects = [
        {
            title: "Projet 1",
            description: "Description détaillée du projet 1, expliquant les technologies utilisées, les défis relevés et les fonctionnalités.",
            image: "assets/images/projet1.jpg",
            link: "#"
        },
        {
            title: "Projet 2",
            description: "Description détaillée du projet 2, expliquant les technologies utilisées, les défis relevés et les fonctionnalités.",
            image: "assets/images/projet2.jpg",
            link: "#"
        },
        {
            title: "Projet 3",
            description: "Description détaillée du projet 3, expliquant les technologies utilisées, les défis relevés et les fonctionnalités.",
            image: "assets/images/projet3.jpg",
            link: "#"
        }
    ];

    const projectList = document.getElementById('project-list');

    projects.forEach(project => {
        const projectCard = document.createElement('div');
        projectCard.className = 'project-card';

        projectCard.innerHTML = `
            <img src="${project.image}" alt="${project.title}">
            <h3>${project.title}</h3>
            <p>${project.description}</p>
            <a href="${project.link}"><button>Voir plus</button></a>
        `;

        projectList.appendChild(projectCard);
    });
});
