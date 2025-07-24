// ==================== SHARED UTILITIES ====================
const Utils = {
    // Initialize theme toggle functionality
    initTheme: function() {
        const themeToggle = document.getElementById('theme-toggle');
        if (!themeToggle) return;

        const themeIcon = themeToggle.querySelector('.material-icons');

        // Set initial theme from localStorage or prefer-color-scheme
        const storedTheme = localStorage.getItem('theme');
        const systemPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        const initialTheme = storedTheme || (systemPrefersDark ? 'dark' : 'light');
        
        document.documentElement.setAttribute('data-theme', initialTheme);
        document.documentElement.setAttribute('data-bs-theme', initialTheme);
        this.updateThemeIcon(initialTheme);

        // Set up click handler
        themeToggle.addEventListener('click', () => {
            const currentTheme = document.documentElement.getAttribute('data-theme');
            const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
            
            // Apply new theme
            document.documentElement.setAttribute('data-theme', newTheme);
            document.documentElement.setAttribute('data-bs-theme', newTheme);
            
            // Store preference
            try {
                localStorage.setItem('theme', newTheme);
            } catch (e) {
                console.error("LocalStorage access denied:", e);
            }
            
            this.updateThemeIcon(newTheme);
        });
    },

    // Update the theme icon
    updateThemeIcon: function(theme) {
        const themeToggle = document.getElementById('theme-toggle');
        if (!themeToggle) return;
        
        const icon = themeToggle.querySelector('.material-icons');
        if (icon) {
            icon.textContent = theme === 'dark' ? 'light_mode' : 'dark_mode';
        }
    },

    // Date formatting
    formatLocalTime: function() {
        document.querySelectorAll(".js-localtime").forEach(el => {
            const raw = el.getAttribute("datetime");
            if (!raw) return;
            el.textContent = new Date(raw).toLocaleString(undefined, { 
                dateStyle: "medium", 
                timeStyle: "short" 
            });
        });
    },

    // Search functionality
initSearch: function(inputId, cardSelector) {
    const searchInput = document.getElementById(inputId);
    // Get a reference to the placeholder
    const noResultsPlaceholder = document.getElementById('noResultsPlaceholder');
    
    if (!searchInput) return;

    searchInput.addEventListener('input', (e) => {
        const searchTerm = e.target.value.toLowerCase();
        const taskCols = document.querySelectorAll(cardSelector);
        let visibleCount = 0; // Counter for visible tasks

        taskCols.forEach(col => {
            const card = col.querySelector('.task-card');
            if (!card) return;
            
            const title = card.querySelector('.card-title')?.textContent.toLowerCase() || '';
            const description = card.querySelector('.card-text')?.textContent.toLowerCase() || '';
            const isVisible = title.includes(searchTerm) || description.includes(searchTerm);
            
            col.style.display = isVisible ? '' : 'none';
            if (isVisible) {
                visibleCount++; // Increment counter if task is shown
            }
        });

        // Show or hide the placeholder based on the count
        if (noResultsPlaceholder) {
            if (visibleCount === 0) {
                noResultsPlaceholder.classList.remove('d-none');
            } else {
                noResultsPlaceholder.classList.add('d-none');
            }
        }
    });

    // Type-to-search functionality (remains the same)
    document.addEventListener('keydown', (e) => {
        if (document.activeElement.tagName === 'INPUT' || 
            document.activeElement.tagName === 'TEXTAREA' || 
            document.querySelector('.modal.show')) {
            return;
        }
        if (e.key.length === 1 && !e.ctrlKey && !e.metaKey) {
            searchInput.focus();
        }
    });
},
    
    // Scramble animation for text
    scrambleAnimate: function(element, finalValue, options = {}) {
        const defaults = { isNumeric: false, charSet: null };
        const config = { ...defaults, ...options };

        let chars;
        if (config.charSet) {
            chars = config.charSet; // Use custom character set if provided
        } else {
            chars = config.isNumeric ? '0123456789' : 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
        }

        let interval = null;
        let counter = 0;
        const duration = 15;
        const speed = 60;

        interval = setInterval(() => {
            element.textContent = String(finalValue).split('').map(char => {
                if (char === ' ' || char === '-') return char;
                return chars[Math.floor(Math.random() * chars.length)];
            }).join('');

            counter++;
            if(counter > duration) {
                clearInterval(interval);
                element.textContent = finalValue;
            }
        }, speed);
    }
};

// ==================== PAGE-SPECIFIC MODULES ====================
const PageModules = {
    // Tasks Page
    tasks: function() {
        // Initialize flatpickr for datetime inputs
        flatpickr(".flatpickr-datetime", {
            enableTime: true,
            dateFormat: "Y-m-d\\TH:i",
            altInput: true,
            altFormat: "F j, Y at h:i K",
            disableMobile: true,
            onChange: function(selectedDates, dateStr, instance) {
                instance.altInput.classList.remove('is-invalid');
            }
        });

        // Form validation
        const addTaskForm = document.getElementById("addTaskForm");
        if (addTaskForm) {
            addTaskForm.addEventListener("submit", function(event) {
                const fpInstance = document.querySelector(".flatpickr-datetime")._flatpickr;
                if (fpInstance.selectedDates.length === 0) {
                    event.preventDefault();
                    event.stopPropagation();
                    fpInstance.altInput.classList.add("is-invalid");
                }
            });
        }

        // Task status highlighting
        const cards = document.querySelectorAll(".task-card");
        const now = new Date();
        cards.forEach(card => {
            const dueRaw = card.dataset.due;
            const isCompleted = card.dataset.completed === "true";
            if (!dueRaw || isCompleted) return;
            
            const due = new Date(dueRaw);
            const msg = card.querySelector(".urgency-message");
            if (!msg) return;
            
            if (now > due) {
                card.classList.add("status-overdue");
                msg.textContent = "Overdue — get it done ASAP.";
            } else if (now.toDateString() === due.toDateString()) {
                card.classList.add("status-due-today");
                msg.textContent = "Due today — complete it while you still can.";
            } else {
                card.classList.add("status-ongoing");
            }
        });

        // Toast notifications
        const toastParam = new URLSearchParams(window.location.search).get("toast");
        if (toastParam) {
            const toastEl = document.getElementById("taskToast");
            const toastMsg = document.getElementById("toastMessage");
            const toastIcon = document.getElementById("toastIcon");

            let message = '';
            let icon = '';
            let stateClass = '';

            switch (toastParam) {
                case "added":
                    message = "Task added successfully.";
                    icon = "add_circle";
                    stateClass = "toast-primary";
                    break;
                case "completed":
                    message = "Task marked as completed.";
                    icon = "check_circle";
                    stateClass = "toast-success";
                    break;
                case "deleted":
                    message = "Task deleted.";
                    icon = "delete";
                    stateClass = "toast-danger";
                    break;
                case "edited":
                    message = "Task description updated.";
                    icon = "edit";
                    stateClass = "toast-warning";
                    break;
                default:
                    return;
            }

            toastMsg.textContent = message;
            toastIcon.textContent = icon;
            toastEl.className = `toast custom-toast ${stateClass}`;
            new bootstrap.Toast(toastEl).show();

            // Clean the URL
            const url = new URL(window.location);
            url.searchParams.delete("toast");
            window.history.replaceState({}, document.title, url);
        }
    },

    // Completed Tasks Page
    completed: function() {
        // Initialize month picker
        flatpickr(".flatpickr-month", {
            plugins: [new monthSelectPlugin({ 
                shorthand: true, 
                dateFormat: "Y-m", 
                altFormat: "F Y", 
                disableMobile: true 
            })]
        });

        // Form validation
        const filterForm = document.getElementById("filterForm");
        if (filterForm) {
            filterForm.addEventListener("submit", function(event) {
                const fpInstance = document.querySelector(".flatpickr-month")._flatpickr;
                if (fpInstance.selectedDates.length === 0) {
                    event.preventDefault();
                    event.stopPropagation();
                    fpInstance.altInput.classList.add("is-invalid");
                }
            });
        }

        // Task status highlighting
        const cards = document.querySelectorAll(".task-card");
        cards.forEach(card => {
            const dueRaw = card.dataset.due;
            const completedRaw = card.dataset.completedAt;
            if (!dueRaw || !completedRaw) return;
            
            const due = new Date(dueRaw);
            const completedAt = new Date(completedRaw);
            const warning = card.querySelector(".late-warning");
            const sameDay = due.toDateString() === completedAt.toDateString();
            
            if (completedAt > due) {
                card.classList.add("overdue");
                if (warning) warning.innerHTML = `<p class="text-danger fw-semibold mb-1">Completed late. Let's aim for green or purple on the next one.</p>`;
            } else if (sameDay) {
                card.classList.add("justintime");
                if (warning) warning.innerHTML = `<p class="text-success fw-semibold mb-1">Just in time, good work.</p>`;
            } else {
                card.classList.add("ahead");
                if (warning) warning.innerHTML = `<p class="text-purple fw-semibold mb-1">Ahead of time, keep it up.</p>`;
            }
        });
    },

    // Overview Page
    overview: function() {
        // Stat boxes animation and data
        const statBoxes = [
            document.getElementById('stat-box-0'),
            document.getElementById('stat-box-1'),
            document.getElementById('stat-box-2'),
            document.getElementById('stat-box-3'),
        ];
        const viewStatusTag = document.getElementById('viewStatusTag');
        let currentView = 'ongoing';
        let currentSort = 'urgency';

        async function updateStatboxes() {
            const response = await fetch(`/api/overview/stats?view=${currentView}&sort=${currentSort}`);
            const data = await response.json();

            statBoxes.forEach((box, i) => {
                const stat = data[i];
                const labelEl = box.querySelector('.stat-label');
                const numberEl = box.querySelector('.stat-number');
                const iconEl = box.querySelector('.icon-placeholder');

                // Update classes for seamless color transition
                box.classList.remove(
                    'border-indigo', 'border-purple', 'border-yellow', 'border-red', 'border-blue', 'border-green',
                    'text-indigo', 'text-purple', 'text-yellow', 'text-red', 'text-blue', 'text-green',
                    'bg-gradient-indigo', 'bg-gradient-purple', 'bg-gradient-yellow', 'bg-gradient-red', 'bg-gradient-blue', 'bg-gradient-green'
                );
                
                // Add new color classes
                box.classList.add(`border-${stat.color}`, `text-${stat.color}`, `bg-gradient-${stat.color}`);
                
                // For numbers (digits only)
                Utils.scrambleAnimate(numberEl, stat.value, {
                    isNumeric: true,
                    duration: 20
                });
                
                // For text (alphanumeric)
                Utils.scrambleAnimate(labelEl, stat.label, {
                    isNumeric: false,
                    duration: 20
                });
                
                // Update icon with smooth transition
                iconEl.textContent = stat.icon;
            });

            // Update status tag
            let statusText = '';
            if (currentView === 'completed') {
                statusText = '<strong>Completed Stats</strong>';
            } else {
                statusText = `<strong>Ongoing by ${currentSort.charAt(0).toUpperCase() + currentSort.slice(1)}</strong>`;
            }
            viewStatusTag.innerHTML = `Viewing: ${statusText}`;
        }

        // Modal Logic
        const applyBtn = document.getElementById('applyViewBtn');
        const viewOptions = document.querySelectorAll('.view-option[data-view]');
        const sortOptions = document.querySelectorAll('.view-option[data-sort]');
        const ongoingSortDiv = document.getElementById('ongoingSortOptions');

        viewOptions.forEach(option => {
            option.addEventListener('click', () => {
                viewOptions.forEach(o => o.classList.remove('active'));
                option.classList.add('active');
                if (option.dataset.view === 'completed') {
                    ongoingSortDiv.classList.add('d-none');
                } else {
                    ongoingSortDiv.classList.remove('d-none');
                }
            });
        });
        
        sortOptions.forEach(option => {
            option.addEventListener('click', () => {
                sortOptions.forEach(o => o.classList.remove('active'));
                option.classList.add('active');
            });
        });

        applyBtn.addEventListener('click', () => {
            currentView = document.querySelector('.view-option[data-view].active').dataset.view;
            if (currentView === 'ongoing') {
                currentSort = document.querySelector('.view-option[data-sort].active').dataset.sort;
            } else {
                currentSort = '';
            }
            updateStatboxes();
        });

        // Initial load with staggered fade-in animation
        statBoxes.forEach((box, index) => {
            setTimeout(() => {
                box.classList.add('fade-in');
            }, index * 150);
        });
        
        setTimeout(() => {
            document.querySelector('.graph-pane').classList.add('fade-in');
        }, 300); 
        
        updateStatboxes();

        // --- Chart Logic ---
        const ctx = document.getElementById('performanceChart');
        let performanceChart; 
        const monthlyBtn = document.getElementById('monthlyBtn');
        const yearlyBtn = document.getElementById('yearlyBtn');
        const chartDescription = document.getElementById('chartDescription');

        async function updateChart(view) {
            const response = await fetch(`/api/overview/chart-data?chart=${view}`);
            const data = await response.json();
            const allScoresAreZero = data.scores.every(score => score === 0);
            if (allScoresAreZero) {
                ctx.classList.add('d-none');
                document.getElementById('chartPlaceholder').classList.remove('d-none');
                if (performanceChart) { performanceChart.destroy(); performanceChart = null; }
                return;
            } else {
                ctx.classList.remove('d-none');
                document.getElementById('chartPlaceholder').classList.add('d-none');
            }
            const backgroundColors = data.scores.map(score => {
                if (score > 1.5) return 'rgba(102, 51, 153, 0.7)';
                if (score < 0.5) return 'rgba(220, 53, 69, 0.7)';
                return 'rgba(25, 135, 84, 0.7)';
            });
            const borderColors = backgroundColors.map(color => color.replace('0.7', '1'));
            if (performanceChart) {
                performanceChart.data.labels = data.labels;
                performanceChart.data.datasets[0].data = data.scores;
                performanceChart.data.datasets[0].backgroundColor = backgroundColors;
                performanceChart.data.datasets[0].borderColor = borderColors;
                performanceChart.update();
            } else {
                performanceChart = new Chart(ctx, {
                    type: 'bar',
                    data: { labels: data.labels, datasets: [{ label: 'Avg. Performance Score', data: data.scores, backgroundColor: backgroundColors, borderColor: borderColors, borderWidth: 1 }] },
                    options: {
                        responsive: true, maintainAspectRatio: false,
                        scales: { y: { beginAtZero: false, ticks: { color: '#6c757d' }, grid: { color: 'rgba(255, 255, 255, 0.1)' } }, x: { ticks: { color: '#6c757d' }, grid: { display: false } } },
                        plugins: { legend: { display: false } }
                    }
                });
            }
        }
        
        monthlyBtn.addEventListener('click', () => {
            updateChart('monthly');
            chartDescription.textContent = 'This graph shows your average performance score for each of the last 12 months.';
            monthlyBtn.querySelector('.material-symbols-outlined').textContent = 'radio_button_checked';
            yearlyBtn.querySelector('.material-symbols-outlined').textContent = 'radio_button_unchecked';
        });
        
        yearlyBtn.addEventListener('click', () => {
            updateChart('yearly');
            chartDescription.textContent = 'This graph shows your average performance score for each of the last 5 years.';
            yearlyBtn.querySelector('.material-symbols-outlined').textContent = 'radio_button_checked';
            monthlyBtn.querySelector('.material-symbols-outlined').textContent = 'radio_button_unchecked';
        });
        
        updateChart('monthly');
    },

    // About Page
    about: function() {
        // Version number scramble animation
        const versionNumber = document.getElementById('version-number');
        if (versionNumber) {
            setTimeout(() => {
                Utils.scrambleAnimate(versionNumber, versionNumber.textContent, {
                    charSet: '0123456789.-abcdefghijklmnopqrstuvwxyz'
                    // charSet: '0123456789.-;'
                });
            }, 30);
        }
    }
};

// ==================== PAGE INITIALIZATION ====================
document.addEventListener("DOMContentLoaded", function() {
    // Initialize shared utilities
    Utils.initTheme();
    Utils.formatLocalTime();

    // Determine which page we're on and initialize specific modules
    const path = window.location.pathname;
    
    if (path === '/' || path === '/overview') {
        Utils.initSearch('taskSearch', '.task-col');
        PageModules.overview();
    } 
    else if (path === '/tasks') {
        Utils.initSearch('taskSearch', '.task-col');
        PageModules.tasks();
    } 
    else if (path === '/completed') {
        Utils.initSearch('taskSearch', '.task-col');
        PageModules.completed();
    } 
    else if (path === '/about') {
        PageModules.about();
    }
});