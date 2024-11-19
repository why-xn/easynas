import React, { useState, useEffect } from "react";
import axios from "axios";
import {useNavigate, useParams} from "react-router-dom";
import constants from "../constants";

const Snapshot = () => {
    const { datasetId } = useParams(); // Get dataset ID from URL
    const { API_URL } = constants;
    const [datasetDetails, setDatasetDetails] = useState(null);
    const [snapshots, setSnapshots] = useState([]);
    const navigate = useNavigate(); // Hook to handle navigation

    const fetchDataset = async () => {
        const authToken = localStorage.getItem("auth_token");
        if (!authToken) {
            navigate('/login'); // Redirect to login if no token is found
            return;
        }

        // Load dataset details
        const datasetRes = await axios.get(`${API_URL}/api/v1/nas/pools/naspool/datasets/${datasetId}`, {
            headers: { Authorization: `${authToken}` },
        });
        setDatasetDetails(datasetRes.data.data);
    }

    const fetchSnapshots = async () => {
        try {
            const authToken = localStorage.getItem("auth_token");
            if (!authToken) {
                navigate('/login'); // Redirect to login if no token is found
                return;
            }
            const response = await axios.get(`${API_URL}/api/v1/nas/pools/naspool/datasets/${datasetId}/snapshots`, {
                headers: {
                    Authorization: `${authToken}`,
                },
            });
            setSnapshots(response.data.data);
        } catch (error) {
            console.error("Error fetching snapshots:", error);
        }
    };

    // Fetch snapshots on page load
    useEffect(() => {
        fetchDataset();
        fetchSnapshots();
    }, [datasetId]);

    // Handle snapshot restore
    const handleRestore = async (snapshotName) => {
        const confirmRestore = window.confirm(
            "Are you sure you want to restore this snapshot?"
        );
        if (!confirmRestore) return;

        try {
            const authToken = localStorage.getItem("auth_token");
            await axios.post(
                `${API_URL}/api/v1/nas/pools/naspool/datasets/${datasetId}/snapshots/restore`,
                {
                    "snapshotName": snapshotName
                },
                {
                    headers: {
                        Authorization: `${authToken}`,
                    },
                }
            );

            alert("Snapshot restored successfully.");
        } catch (error) {
            console.error("Error restoring snapshot:", error);
            alert("Failed to restore snapshot.");
        }
    };

    // Handle snapshot delete
    const handleDelete = async (snapshotName) => {
        const confirmDelete = window.confirm(
            "Are you sure you want to delete this snapshot?"
        );
        if (!confirmDelete) return;

        try {
            // Convert the path to Base64 encoding
            const encodedSnapshotName = btoa(snapshotName);

            const authToken = localStorage.getItem("auth_token");
            await axios.delete(`${API_URL}/api/v1/nas/pools/naspool/datasets/${datasetId}/snapshots/${encodedSnapshotName}`, {
                headers: {
                    Authorization: `${authToken}`,
                },
            });

            setSnapshots(snapshots.filter((snapshot) => snapshot.name !== snapshotName));
            alert("Snapshot deleted successfully.");
        } catch (error) {
            console.error("Error deleting snapshot:", error);
            alert("Failed to delete snapshot.");
        }
    };

    // Handle add snapshot
    const handleAddSnapshot = async () => {
        const confirmAdd = window.confirm(
            "Are you sure you want to create a new snapshot?"
        );
        if (!confirmAdd) return;

        try {
            const authToken = localStorage.getItem("auth_token");
            const response = await axios.post(
                `${API_URL}/api/v1/nas/pools/naspool/datasets/${datasetId}/snapshots`,
                {},
                {
                    headers: {
                        Authorization: `${authToken}`,
                    },
                }
            );

            fetchSnapshots()
            alert("Snapshot created successfully.");
        } catch (error) {
            console.error("Error creating snapshot:", error);
            alert("Failed to create snapshot.");
        }
    };

    return (
        <div style={{ padding: "20px" }}>
            <h1 style={{ marginBottom: "20px" }}>Snapshots</h1>

            {/* Add Snapshot Button */}
            <button
                onClick={handleAddSnapshot}
                style={{
                    marginBottom: "20px",
                    padding: "10px 20px",
                    backgroundColor: "#007bff",
                    color: "#fff",
                    border: "none",
                    borderRadius: "5px",
                    cursor: "pointer",
                }}
            >
                Add Snapshot
            </button>

            {/* Snapshot Table */}
            <table
                style={{
                    width: "100%",
                    borderCollapse: "collapse",
                }}
            >
                <thead>
                <tr
                    style={{
                        backgroundColor: "#f4f4f4",
                        borderBottom: "2px solid #ddd",
                        textAlign: "left",
                    }}
                >
                    <th style={{ padding: "12px" }}>Snapshot Name</th>
                    <th style={{ padding: "12px" }}>User</th>
                    <th style={{ padding: "12px" }}>Referenced</th>
                    <th style={{ padding: "12px" }}>Created At</th>
                    <th style={{ padding: "12px", textAlign: "center" }}>Actions</th>
                </tr>
                </thead>
                <tbody>
                {snapshots?.map((snapshot) => (
                    <tr
                        key={snapshot.name}
                        style={{
                            borderBottom: "1px solid #ddd",
                        }}
                    >
                        <td style={{ padding: "12px" }}>{snapshot.name}</td>
                        <td style={{ padding: "12px" }}>{snapshot.used}</td>
                        <td style={{ padding: "12px" }}>{snapshot.referenced}</td>
                        <td style={{ padding: "12px" }}>{snapshot.createdAt}</td>
                        <td
                            style={{
                                padding: "12px",
                                textAlign: "center",
                            }}
                        >
                            <button
                                onClick={() => handleRestore(snapshot.name)}
                                style={{
                                    backgroundColor: "#28a745",
                                    color: "#fff",
                                    border: "none",
                                    borderRadius: "5px",
                                    cursor: "pointer",
                                    padding: "5px 10px",
                                    marginRight: "10px",
                                }}
                            >
                                Restore
                            </button>
                            <button
                                onClick={() => handleDelete(snapshot.name)}
                                style={{
                                    backgroundColor: "#dc3545",
                                    color: "#fff",
                                    border: "none",
                                    borderRadius: "5px",
                                    cursor: "pointer",
                                    padding: "5px 10px",
                                }}
                            >
                                Delete
                            </button>
                        </td>
                    </tr>
                ))}
                </tbody>
            </table>
        </div>
    );
};

export default Snapshot;
