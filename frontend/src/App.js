import React from 'react';
import {BrowserRouter as Router, Navigate, Route, Routes, useLocation} from 'react-router-dom';
import Dashboard from './components/Dashboard';
import Users from './components/Users';
import Sidenav from './components/Sidenav';
import './App.css';
import Login from "./components/Login";
import DatasetDetails from "./components/DatasetDetails";
import Filesystem from "./components/Filesystem";
import Snapshots from "./components/Snapshots"; // Import global CSS

const App = () => {
    const location = useLocation(); // Get the current location
    // Conditionally add class for login page to remove margins
    const isLoginPage = location.pathname === '/login';

    // Function to check if the token is valid
    const isAuthenticated = () => {
        const authToken = localStorage.getItem("auth_token");
        if (!authToken) return false;

        // You can add more checks here if your token has an expiration date, etc.
        return true;
    };

    return (
        <div className="app-container">
            {location.pathname !== '/login' && <Sidenav />}
            <div className={`main-content ${isLoginPage ? 'login-page' : ''}`}>
                <Routes>
                    {/* Redirect to Login page if no token */}
                    <Route path="/login" element={<Login />} />

                    {/* Protected Routes */}
                    <Route path="/" element={isAuthenticated() ? <Dashboard /> : <Navigate to="/login" />} />
                    <Route path="/dashboard" element={isAuthenticated() ? <Dashboard /> : <Navigate to="/login" />} />
                    <Route path="/users" element={isAuthenticated() ? <Users /> : <Navigate to="/login" />} />
                    <Route path="/dataset/:id" element={isAuthenticated() ? <DatasetDetails /> : <Navigate to="/login" />} />
                    <Route path="/dataset/:datasetId/filesystem" element={isAuthenticated() ? <Filesystem /> : <Navigate to="/login" />} />
                    <Route path="/dataset/:datasetId/snapshots" element={isAuthenticated() ? <Snapshots /> : <Navigate to="/login" />} />
                </Routes>
            </div>
        </div>
    );
};

export default App;
