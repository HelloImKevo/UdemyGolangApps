// Main JavaScript functionality for the Login App

// Utility functions
const utils = {
    // Show notification
    showNotification: function(message, type = 'info') {
        const notification = document.createElement('div');
        notification.className = `notification ${type}`;
        notification.textContent = message;
        
        document.body.appendChild(notification);
        
        // Auto remove after 5 seconds
        setTimeout(() => {
            if (notification.parentNode) {
                notification.parentNode.removeChild(notification);
            }
        }, 5000);
    },

    // Format date for display
    formatDate: function(dateString) {
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        });
    },

    // Validate email format
    isValidEmail: function(email) {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailRegex.test(email);
    },

    // Get token from localStorage
    getAuthToken: function() {
        return localStorage.getItem('authToken');
    },

    // Remove auth data
    clearAuth: function() {
        localStorage.removeItem('authToken');
        localStorage.removeItem('user');
    },

    // Check if user is authenticated
    isAuthenticated: function() {
        return !!this.getAuthToken();
    }
};

// API helper functions
const api = {
    // Base API call
    call: async function(endpoint, options = {}) {
        const defaultOptions = {
            headers: {
                'Content-Type': 'application/json',
            }
        };

        // Add auth token if available
        const token = utils.getAuthToken();
        if (token) {
            defaultOptions.headers.Authorization = `Bearer ${token}`;
        }

        const config = {
            ...defaultOptions,
            ...options,
            headers: {
                ...defaultOptions.headers,
                ...options.headers
            }
        };

        try {
            const response = await fetch(endpoint, config);
            const data = await response.json();
            
            return {
                success: response.ok,
                status: response.status,
                data: data
            };
        } catch (error) {
            return {
                success: false,
                status: 0,
                error: error.message
            };
        }
    },

    // Login
    login: async function(credentials) {
        return this.call('/api/auth/login', {
            method: 'POST',
            body: JSON.stringify(credentials)
        });
    },

    // Register
    register: async function(userData) {
        return this.call('/api/auth/register', {
            method: 'POST',
            body: JSON.stringify(userData)
        });
    },

    // Get profile
    getProfile: async function() {
        return this.call('/api/auth/profile', {
            method: 'GET'
        });
    },

    // Logout
    logout: async function() {
        return this.call('/api/auth/logout', {
            method: 'POST'
        });
    }
};

// Navigation functions
const navigation = {
    // Update navigation based on auth status
    updateNavigation: function() {
        const isAuth = utils.isAuthenticated();
        const navMenu = document.querySelector('.nav-menu');
        
        if (!navMenu) return;

        if (isAuth) {
            // Show authenticated navigation
            navMenu.innerHTML = `
                <li class="nav-item">
                    <a href="/dashboard" class="nav-link">Dashboard</a>
                </li>
                <li class="nav-item">
                    <a href="#" class="nav-link" onclick="navigation.logout()">Logout</a>
                </li>
            `;
        } else {
            // Show guest navigation
            navMenu.innerHTML = `
                <li class="nav-item">
                    <a href="/" class="nav-link">Home</a>
                </li>
                <li class="nav-item">
                    <a href="/login" class="nav-link">Login</a>
                </li>
                <li class="nav-item">
                    <a href="/register" class="nav-link">Register</a>
                </li>
            `;
        }
    },

    // Logout function
    logout: async function() {
        try {
            await api.logout();
            utils.clearAuth();
            utils.showNotification('Logged out successfully', 'success');
            window.location.href = '/';
        } catch (error) {
            utils.clearAuth();
            window.location.href = '/';
        }
    }
};

// Form validation
const validation = {
    // Validate login form
    validateLoginForm: function(formData) {
        const errors = [];

        if (!formData.email) {
            errors.push('Email is required');
        } else if (!utils.isValidEmail(formData.email)) {
            errors.push('Please enter a valid email address');
        }

        if (!formData.password) {
            errors.push('Password is required');
        } else if (formData.password.length < 6) {
            errors.push('Password must be at least 6 characters long');
        }

        return errors;
    },

    // Validate registration form
    validateRegisterForm: function(formData) {
        const errors = [];

        if (!formData.first_name || formData.first_name.trim().length < 1) {
            errors.push('First name is required');
        }

        if (!formData.last_name || formData.last_name.trim().length < 1) {
            errors.push('Last name is required');
        }

        if (!formData.username || formData.username.length < 3) {
            errors.push('Username must be at least 3 characters long');
        }

        if (!formData.email) {
            errors.push('Email is required');
        } else if (!utils.isValidEmail(formData.email)) {
            errors.push('Please enter a valid email address');
        }

        if (!formData.password) {
            errors.push('Password is required');
        } else if (formData.password.length < 6) {
            errors.push('Password must be at least 6 characters long');
        }

        return errors;
    }
};

// Initialize app when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    // Update navigation
    navigation.updateNavigation();

    // Add some interactive enhancements
    const buttons = document.querySelectorAll('.btn');
    buttons.forEach(button => {
        button.addEventListener('click', function(e) {
            // Add click effect
            this.style.transform = 'scale(0.95)';
            setTimeout(() => {
                this.style.transform = '';
            }, 150);
        });
    });

    // Add form enhancements
    const inputs = document.querySelectorAll('input');
    inputs.forEach(input => {
        input.addEventListener('focus', function() {
            this.parentNode.classList.add('focused');
        });

        input.addEventListener('blur', function() {
            this.parentNode.classList.remove('focused');
        });
    });
});

// Check authentication on protected pages
function checkAuthOnProtectedPage() {
    const protectedPaths = ['/dashboard'];
    const currentPath = window.location.pathname;

    if (protectedPaths.includes(currentPath) && !utils.isAuthenticated()) {
        window.location.href = '/login';
    }
}

// Check redirect on auth pages when already authenticated
function checkAuthOnGuestPage() {
    const guestPaths = ['/login', '/register'];
    const currentPath = window.location.pathname;

    if (guestPaths.includes(currentPath) && utils.isAuthenticated()) {
        window.location.href = '/dashboard';
    }
}

// Run auth checks
checkAuthOnProtectedPage();
checkAuthOnGuestPage();

// Export for global access
window.loginApp = {
    utils,
    api,
    navigation,
    validation
};
