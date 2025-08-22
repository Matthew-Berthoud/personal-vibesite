// --- Theme Toggler ---
const themeToggle = document.getElementById('theme-toggle');
const body = document.body;

// Function to set theme
const setTheme = (theme) => {
	body.setAttribute('data-theme', theme);
};

// Event listener for the button
themeToggle.addEventListener('click', () => {
	const currentTheme = body.getAttribute('data-theme');
	if (currentTheme === 'dark') {
		setTheme('light');
	} else {
		setTheme('dark');
	}
});

// --- Smooth Scrolling for Navigation Links ---
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
	anchor.addEventListener('click', function(e) {
		e.preventDefault();
		document.querySelector(this.getAttribute('href')).scrollIntoView({
			behavior: 'smooth'
		});
	});
});

// --- Contact Form Submission ---
// This is a basic example. For a real website, you would need a backend service to handle form submissions.
const contactForm = document.getElementById('contact-form');
contactForm.addEventListener('submit', function(e) {
	e.preventDefault();
	const name = this.querySelector('input[name="name"]').value;
	// A simple confirmation message instead of an alert
	const confirmationMessage = document.createElement('p');
	confirmationMessage.textContent = `Thank you, ${name}! Your message has been "sent". (This is a demo and does not actually send emails).`;
	confirmationMessage.style.textAlign = 'center';
	confirmationMessage.style.marginTop = '20px';
	confirmationMessage.style.color = 'var(--secondary-color)';
	this.insertAdjacentElement('afterend', confirmationMessage);

	this.reset();

	// Remove the message after a few seconds
	setTimeout(() => {
		confirmationMessage.remove();
	}, 5000);
});
