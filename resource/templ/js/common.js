// Common.js - Connect Four Auth Functions
console.log("Connect Four common.js loaded successfully");

// Functions to handle JWT token
function saveToken(token) {
    try {
        // First, test if localStorage is accessible and working
        localStorage.setItem('testKey', 'testValue');
        const testValue = localStorage.getItem('testKey');
        if (testValue !== 'testValue') {
            console.error("LocalStorage test failed: stored value doesn't match retrieved value");
            // Try to fall back to cookies
            return saveTokenCookie(token);
        }
        localStorage.removeItem('testKey');
        
        // Now save the actual token
        console.log("Saving token to localStorage:", token.substring(0, 15) + "...");
        localStorage.setItem('jwtToken', token);
        
        // Verify the token was saved
        const savedToken = localStorage.getItem('jwtToken');
        if (!savedToken) {
            console.error("Failed to save token to localStorage");
            // Try to fall back to cookies
            return saveTokenCookie(token);
        }
        console.log("Token successfully saved to localStorage");
        return true;
    } catch (e) {
        console.error("Error accessing localStorage:", e);
        // Try to fall back to cookies
        return saveTokenCookie(token);
    }
}

// Fallback method using cookies
function saveTokenCookie(token) {
    try {
        console.log("Trying to save token as cookie");
        // Set cookie to expire in 1 day
        const expiration = new Date();
        expiration.setDate(expiration.getDate() + 1);
        
        document.cookie = `jwtToken=${token}; expires=${expiration.toUTCString()}; path=/`;
        
        // Test if cookie was set
        if (document.cookie.indexOf('jwtToken=') >= 0) {
            console.log("Token successfully saved as cookie");
            return true;
        } else {
            console.error("Failed to save token as cookie");
            alert("Failed to save login information. Please check your browser settings.");
            return false;
        }
    } catch (e) {
        console.error("Error setting cookie:", e);
        alert("Error saving login information. This could be due to private browsing mode or disabled cookies.");
        return false;
    }
}

function getToken() {
    try {
        // First try localStorage
        const token = localStorage.getItem('jwtToken');
        if (token) {
            console.log("Token retrieved from localStorage:", token.substring(0, 15) + "...");
            return token;
        }
        
        // If not in localStorage, try cookies
        console.log("Token not found in localStorage, checking cookies");
        const cookies = document.cookie.split(';');
        for (let i = 0; i < cookies.length; i++) {
            const cookie = cookies[i].trim();
            if (cookie.startsWith('jwtToken=')) {
                const cookieToken = cookie.substring('jwtToken='.length);
                console.log("Token retrieved from cookie:", cookieToken.substring(0, 15) + "...");
                return cookieToken;
            }
        }
        
        console.log("No JWT token found in localStorage or cookies");
        return null;
    } catch (e) {
        console.error("Error retrieving token:", e);
        return null;
    }
}

function clearToken() {
    try {
        console.log("Clearing token from storage");
        // Clear from localStorage
        localStorage.removeItem('jwtToken');
        
        // Clear from cookies
        document.cookie = "jwtToken=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    } catch (e) {
        console.error("Error clearing token:", e);
    }
}

function isLoggedIn() {
    try {
        const isLogged = !!getToken();
        console.log("isLoggedIn check:", isLogged);
        return isLogged;
    } catch (e) {
        console.error("Error in isLoggedIn check:", e);
        return false;
    }
}

// Centralized logout function
function logout() {
    clearToken();
    window.location.href = '/ui';
}

// HTMX extension to add JWT auth headers to requests
htmx.defineExtension('jwt-auth', {
    onEvent: function(name, evt) {
        // When content is about to be requested, add the JWT header if available
        if (name === 'htmx:configRequest') {
            const token = getToken();
            if (token) {
                console.log("Adding Authorization header to request:", evt.detail.path);
                evt.detail.headers['Authorization'] = 'Bearer ' + token;
            } else {
                console.log("No token available for request:", evt.detail.path);
            }
        }
    }
});

// Function to format dates (for client-side rendering)
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: 'numeric',
        minute: '2-digit',
        hour12: true
    });
}
