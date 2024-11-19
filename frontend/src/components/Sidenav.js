import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import './Sidenav.css';

const Sidenav = () => {
    const navigate = useNavigate(); // Hook for navigation

    const handleLogout = () => {
        localStorage.removeItem('auth_token'); // Remove the auth token
        navigate('/login'); // Redirect to login page
    };

    return (
        <div className="sidenav">
            <h2>easyNAS</h2>
            <Link to="/dashboard">Dashboard</Link>
            <Link to="/users">Users</Link>

            {/* Logout Button */}
            <button className="logout-button" onClick={handleLogout}>
                Logout
            </button>
        </div>
    );
};

export default Sidenav;
