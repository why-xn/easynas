import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom'; // Use useNavigate for routing
import { Link } from 'react-router-dom'; // Import Link for routing
import { CircularProgressbar } from 'react-circular-progressbar'; // Import CircularProgressbar
import 'react-circular-progressbar/dist/styles.css'; // Import the styles for CircularProgressbar

import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    BarElement,
    LineElement,
    Title,
    Tooltip,
    Legend,
    PointElement, // Make sure PointElement is registered
} from 'chart.js';
import './Dashboard.css';
import constants from "../constants";

// Register the necessary chart components
ChartJS.register(
    CategoryScale,
    LinearScale,
    BarElement,
    LineElement,
    Title,
    Tooltip,
    Legend,
    PointElement // Register PointElement explicitly for line charts
);

const Dashboard = () => {
    const { API_URL } = constants;
    const [metrics, setMetrics] = useState(null);
    const [datasets, setDatasets] = useState([]);
    const [showModal, setShowModal] = useState(false);
    const [newDatasetName, setNewDatasetName] = useState("");
    const [newQuota, setNewQuota] = useState("");
    const navigate = useNavigate(); // Hook to handle navigation

    const token = localStorage.getItem('auth_token');

    const fetchDatasets = async () => {
        try {
            const response = await axios.get(`${API_URL}/api/v1/nas/pools/naspool/datasets`, {
                headers: { Authorization: `${token}` },
            });
            setDatasets(response.data.data);
        } catch (err) {
            console.error("Error fetching datasets", err);
        }
    };

    const fetchMetrics = async () => {
        try {
            const response = await axios.get(`${API_URL}/api/v1/metrics/system`, {
                headers: { Authorization: `${token}` },
            });
            setMetrics(response.data.data);
        } catch (err) {
            console.error("Error fetching system metrics", err);
        }
    };

    useEffect(() => {
        if (!token) {
            navigate('/login'); // Redirect to login if no token is found
            return;
        }
        fetchMetrics();
        fetchDatasets();

    }, []);

    // Handle form submission for creating a dataset
    const handleCreateDataset = async (e) => {
        e.preventDefault();

        try {
            const payload = {
                datasetName: newDatasetName,
                quota: newQuota,
            };

            await axios.post(`${API_URL}/api/v1/nas/pools/naspool/datasets`, payload, {
                headers: { Authorization: `${token}` },
            });

            alert("Dataset created successfully!");

            // Reload datasets
            fetchDatasets();

            // Close modal
            setShowModal(false);
            setNewDatasetName("");
            setNewQuota("");
        } catch (error) {
            console.error("Error creating dataset:", error);
            alert("Failed to create dataset. Please try again.");
        }
    };

    if (!metrics) {
        return <div>Loading...</div>;
    }

    // Helper functions to convert bytes to GB, TB
    const bytesToGB = (bytes) => (bytes / (1024 ** 3)).toFixed(2);
    const bytesToTB = (bytes) => (bytes / (1024 ** 4)).toFixed(2);

    return (
        <div className="dashboard">
            <h2>Dashboard</h2>
            <div className="metric-cards">
                <div className="metric-card">
                    <h3>CPU Usage</h3>
                    <CircularProgressbar
                        value={metrics.cpuUsagePercent}
                        text={`${metrics.cpuUsagePercent.toFixed(1)}%`}
                        strokeWidth={10}
                    />
                    <p>Total CPUs: {metrics.totalCpus}</p>
                </div>
                <div className="metric-card">
                    <h3>Memory Usage</h3>
                    <CircularProgressbar
                        value={metrics.memoryPercent}
                        text={`${metrics.memoryPercent.toFixed(1)}%`}
                        strokeWidth={10}
                    />
                    <p>Total Memory: {bytesToGB(metrics.totalMemory)} GB</p>
                </div>
                <div className="metric-card">
                    <h3>Disk Usage</h3>
                    <CircularProgressbar
                        value={metrics.diskPercent}
                        text={`${metrics.diskPercent.toFixed(1)}%`}
                        strokeWidth={10}
                    />
                    <p>Total Disk: {bytesToGB(metrics.totalDisk)} GB</p>
                </div>
            </div>

            <h2 className="mt-20">ZFS Datasets</h2>
            <button
                style={{
                    padding: "10px 20px",
                    backgroundColor: "#2196f3",
                    color: "#fff",
                    border: "none",
                    borderRadius: "5px",
                    marginTop: "0px",
                    cursor: "pointer",
                    marginBottom: "15px"
                }}
                onClick={() => setShowModal(true)}>
                    Create New Dataset
            </button>
            <table className="zfs-table">
                <thead>
                <tr>
                    <th>Name</th>
                    <th>Quota</th>
                    <th>Used</th>
                    <th>Available</th>
                    <th>Share Enabled</th>
                </tr>
                </thead>
                <tbody>
                {datasets.map((dataset) => (
                    <tr key={dataset.id}>
                        <td>
                            {<Link to={`/dataset/${dataset.id}`}>{dataset.name}</Link>}
                        </td>
                        <td>{dataset.quota}</td>
                        <td>{dataset.used}</td>
                        <td>{dataset.available}</td>
                        <td>{dataset.shareEnabled ? 'Yes' : 'No'}</td>
                    </tr>
                ))}
                </tbody>
            </table>

            {/* Basic HTML Modal */}
            {showModal && (
                <div
                    id="datasetModal"
                    style={{
                        position: "fixed",
                        top: "0",
                        left: "0",
                        width: "100%",
                        height: "100%",
                        backgroundColor: "rgba(0, 0, 0, 0.5)",
                        display: "flex",
                        alignItems: "center",
                        justifyContent: "center",
                    }}
                >
                    <div
                        style={{
                            backgroundColor: "#fff",
                            padding: "20px",
                            borderRadius: "5px",
                            width: "400px",
                        }}
                    >
                        <h3>Create New Dataset</h3>
                        <form onSubmit={handleCreateDataset}>
                            <div style={{ marginBottom: "10px" }}>
                                <label htmlFor="datasetName">Dataset Name:</label>
                                <input
                                    type="text"
                                    id="datasetName"
                                    name="datasetName"
                                    value={newDatasetName}
                                    onChange={(e) => setNewDatasetName(e.target.value)}
                                    required
                                    style={{
                                        width: "100%",
                                        padding: "8px",
                                        marginTop: "5px",
                                        border: "1px solid #ccc",
                                        borderRadius: "4px",
                                    }}
                                />
                            </div>
                            <div style={{ marginBottom: "10px" }}>
                                <label htmlFor="quota">Quota:</label>
                                <input
                                    type="text"
                                    id="quota"
                                    name="quota"
                                    value={newQuota}
                                    onChange={(e) => setNewQuota(e.target.value)}
                                    required
                                    style={{
                                        width: "100%",
                                        padding: "8px",
                                        marginTop: "5px",
                                        border: "1px solid #ccc",
                                        borderRadius: "4px",
                                    }}
                                />
                            </div>
                            <div style={{ display: "flex", justifyContent: "space-between" }}>
                                <button
                                    type="button"
                                    className="btn btn-secondary"
                                    style={{
                                        padding: "10px 20px",
                                        backgroundColor: "#f44336",
                                        color: "#fff",
                                        border: "none",
                                        borderRadius: "5px",
                                        cursor: "pointer",
                                    }}
                                    onClick={() => setShowModal(false)}
                                >
                                    Cancel
                                </button>
                                <button
                                    type="submit"
                                    className="btn btn-primary"
                                    style={{
                                        padding: "10px 20px",
                                        backgroundColor: "#4caf50",
                                        color: "#fff",
                                        border: "none",
                                        borderRadius: "5px",
                                        cursor: "pointer",
                                    }}>
                                    Create Dataset
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

        </div>
    );
};

export default Dashboard;
