// TutorGPT Chat Application
// Handles the chat interface, session management, and API calls

document.addEventListener('DOMContentLoaded', () => {
    // DOM Elements
    const chatForm = document.getElementById('chat-form');
    const userInput = document.getElementById('user-input');
    const chatMessages = document.getElementById('chat-messages');
    
    // API Configuration
    const API_ENDPOINT = 'http://localhost:8000/api/chat';
    
    // Debug Logging
    const DEBUG = false;
    const logs = [];
    
    const debugLog = (message, data) => {
        if (!DEBUG) return;
        
        const timestamp = new Date().toISOString();
        const logEntry = {
            timestamp,
            message,
            data: data ? JSON.stringify(data) : undefined
        };
        
        logs.push(logEntry);
        console.log(`[${timestamp}] ${message}`, data || '');
    };
    
    // Session Management
    let sessionHistory = [];
    
    // Load session (only initializes a new session)
    const loadSession = () => {
        // Always start with empty session
        sessionHistory = [];
        debugLog('Session initialized with empty history');
        
        // Add a system message about session persistence
        addMessageToUI('system', 'Chat history is not saved and will be cleared when you refresh or close this page.');
        
        // Send welcome message for new session
        debugLog('New session detected, requesting welcome message');
        sendMessage("", "welcome", false);
    };
    
    // Save session - no need to persist
    const saveSession = () => {
        // No persistence, just keep in memory
        debugLog('Session updated in memory', { historyLength: sessionHistory.length });
    };
    
    // Format message content with basic markdown-like styling
    const formatMessage = (content) => {
        // Replace code blocks with formatted code
        content = content.replace(/```([^`]+)```/g, '<pre><code>$1</code></pre>');
        
        // Replace inline code with formatted code
        content = content.replace(/`([^`]+)`/g, '<code>$1</code>');
        
        // Replace line breaks with <br>
        content = content.replace(/\n/g, '<br>');
        
        // Wrap in paragraph
        return `<p>${content}</p>`;
    };
    
    // Render visualization content
    const renderVisualization = (payload, parentElement) => {
        const visualizationContainer = document.createElement('div');
        visualizationContainer.className = 'visualization-container';
        
        // Check if payload is for an iframe (has HTML doctype or <html> tag)
        if (payload.includes('<!DOCTYPE html>') || payload.includes('<html')) {
            // Create a sandbox iframe for full HTML documents
            const iframe = document.createElement('iframe');
            iframe.sandbox = 'allow-scripts allow-same-origin allow-popups allow-forms';
            iframe.loading = 'lazy';
            
            // Append iframe to container first
            visualizationContainer.appendChild(iframe);
            
            // Set content after iframe is in the DOM
            setTimeout(() => {
                const iframeDoc = iframe.contentDocument || iframe.contentWindow.document;
                iframeDoc.open();
                iframeDoc.write(payload);
                iframeDoc.close();
            }, 0);
        } else {
            // For regular HTML fragments
            visualizationContainer.innerHTML = payload;
            
            // Execute any scripts in the visualization (if scripts are included)
            const scripts = visualizationContainer.querySelectorAll('script');
            scripts.forEach(script => {
                const newScript = document.createElement('script');
                
                // Copy all attributes
                Array.from(script.attributes).forEach(attr => {
                    newScript.setAttribute(attr.name, attr.value);
                });
                
                // Copy inline script content
                newScript.textContent = script.textContent;
                
                // Replace the old script with the new one
                script.parentNode.replaceChild(newScript, script);
            });
        }
        
        // If parent element is provided, append to it, otherwise append to the message container
        if (parentElement) {
            parentElement.appendChild(visualizationContainer);
        } else {
            // Find the last message and append to it
            const messages = document.querySelectorAll('.message');
            if (messages.length > 0) {
                messages[messages.length - 1].appendChild(visualizationContainer);
            }
        }
    };
    
    // Add a message to the UI
    const addMessageToUI = (role, content, visualizationPayload = null) => {
        const messageDiv = document.createElement('div');
        messageDiv.className = `message ${role}`;
        
        // Create the message content div
        const contentDiv = document.createElement('div');
        contentDiv.className = 'message-content';
        
        if (role === 'error') {
            // For error messages, we want to style them differently
            contentDiv.innerHTML = `<p>${content}</p>`;
        } else {
            // For normal messages, we'll support basic markdown-like formatting
            contentDiv.innerHTML = formatMessage(content);
            
            // Add to session history for user and assistant messages
            if (role === 'user' || role === 'assistant') {
                sessionHistory.push({ role, content });
                saveSession();
            }
        }
        
        messageDiv.appendChild(contentDiv);
        
        // Add visualization if present
        if (visualizationPayload) {
            renderVisualization(visualizationPayload, messageDiv);
        }
        
        // Add to message container and scroll to view
        chatMessages.appendChild(messageDiv);
        chatMessages.scrollTop = chatMessages.scrollHeight;
    };
    
    // Show typing indicator
    const showTypingIndicator = () => {
        debugLog('Showing typing indicator');
        
        // Check if the indicator already exists (to avoid duplicates)
        if (document.getElementById('typing-indicator')) {
            debugLog('Typing indicator already exists, not creating a new one');
            return;
        }
        
        const typingDiv = document.createElement('div');
        typingDiv.className = 'message assistant';
        typingDiv.id = 'typing-indicator';
        
        typingDiv.innerHTML = `
            <div class="typing-indicator">
                <span></span>
                <span></span>
                <span></span>
            </div>
        `;
        
        chatMessages.appendChild(typingDiv);
        chatMessages.scrollTop = chatMessages.scrollHeight;
    };
    
    // Hide typing indicator
    const hideTypingIndicator = () => {
        debugLog('Hiding typing indicator');
        const typingIndicator = document.getElementById('typing-indicator');
        if (typingIndicator) {
            typingIndicator.remove();
        } else {
            debugLog('WARNING: Tried to hide typing indicator but it was not found in the DOM');
        }
    };
    
    // Send message to the backend API
    const sendMessage = async (message, messageType = 'chat', addToUI = true) => {
        try {
            if (addToUI && message.trim() !== '') {
                // Add user message to UI
                addMessageToUI('user', message);
            }
            
            // Show the typing indicator
            showTypingIndicator();
            
            // Set the global flag to indicate API request is in progress
            window.apiRequestInProgress = true;
            debugLog(`Sending ${messageType} message to API`, { message, historyLength: sessionHistory.length });
            
            let response;
            try {
                response = await fetch(API_ENDPOINT, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        session_history: sessionHistory,
                        message_type: messageType,
                        content: message
                    }),
                    cache: 'no-store'
                });
                
                debugLog('Fetch request completed', { 
                    status: response.status, 
                    ok: response.ok
                });
            } catch (fetchError) {
                debugLog('Network error when fetching', { error: fetchError.message });
                hideTypingIndicator();
                addMessageToUI('error', 'Network error. Please check your connection and try again.');
                window.apiRequestInProgress = false;
                return;
            }
            
            if (!response.ok) {
                debugLog('Response not OK', { status: response.status });
                hideTypingIndicator();
                addMessageToUI('error', `Server error (${response.status}). Please try again later.`);
                window.apiRequestInProgress = false;
                return;
            }
            
            // Process the response
            const data = await response.json();
            debugLog('Received API response', { 
                hasResponse: !!data.response, 
                hasError: !!data.error,
                hasVisualization: !!data.visualization_payload
            });
            
            // Hide typing indicator
            hideTypingIndicator();
            
            if (data.error) {
                debugLog('Response contains error field', { error: data.error });
                
                // If we have a user-friendly message in the response, use that
                if (data.response) {
                    addMessageToUI('assistant', data.response);
                } else {
                    addMessageToUI('error', `Error: ${data.error}`);
                }
            } else if (data.response) {
                // Add the response to the UI
                addMessageToUI('assistant', data.response);
                
                // If there's a visualization payload, render it
                if (data.visualization_payload) {
                    debugLog('Rendering visualization payload');
                    renderVisualization(data.visualization_payload);
                }
            } else {
                // We got valid JSON but no response field
                debugLog('Response field missing in API response');
                addMessageToUI('error', 'Received an incomplete response from the server.');
            }
        } catch (error) {
            debugLog('Error in sendMessage', { error: error.message });
            hideTypingIndicator();
            addMessageToUI('error', 'Something went wrong. Please try again later.');
        } finally {
            // Clear the API request flag
            window.apiRequestInProgress = false;
            debugLog('API request completed and flag cleared');
        }
    };
    
    // Handle form submission
    chatForm.addEventListener('submit', (e) => {
        e.preventDefault();
        
        const message = userInput.value.trim();
        if (!message) return;
        
        // Determine message type (chat or visualization)
        let messageType = 'chat';
        
        // Keywords that suggest a visualization request
        const visualizationKeywords = [
            'visualize', 'visualization', 'visualizer', 
            'show me', 'diagram', 'graph', 'chart', 'plot',
            'interactive', 'animation', 'animate', 'visual',
            'display', 'illustrate', 'demonstrate', 'simulation'
        ];
        
        // Check for visualization keywords
        const lowerMessage = message.toLowerCase();
        if (visualizationKeywords.some(keyword => lowerMessage.includes(keyword))) {
            messageType = 'visualization';
        } 
        // Check for task keywords
        else if (lowerMessage.includes('task') ||
                 lowerMessage.includes('exercise') ||
                 lowerMessage.includes('problem') ||
                 lowerMessage.includes('challenge')) {
            messageType = 'task';
        }
        
        // Send message to backend (true = add to UI)
        sendMessage(message, messageType, true);
        
        // Clear input
        userInput.value = '';
    });
    
    // Handle Shift+Enter for new lines
    userInput.addEventListener('keydown', (e) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            chatForm.dispatchEvent(new Event('submit'));
        }
    });
    
    // Initialize
    loadSession();
    
    // Focus on input
    userInput.focus();
}); 