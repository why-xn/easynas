import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router } from 'react-router-dom'; // Import BrowserRouter
import App from './App'; // Import your main App component
import './index.css'; // Global styles

ReactDOM.render(
    <Router> {/* Wrap App with Router */}
        <App />
    </Router>,
    document.getElementById('root')
);
