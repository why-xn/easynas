import React, { useState, useEffect } from "react";
import axios from "axios";
import {useNavigate, useParams} from "react-router-dom";
import constants from "../constants";

const Filesystem = () => {
    const { datasetId } = useParams(); // Get dataset ID from URL
    const { API_URL } = constants;
    const [datasetDetails, setDatasetDetails] = useState(null);
    const [filesystem, setFilesystem] = useState([]);
    const [viewMode, setViewMode] = useState("list"); // "list" or "grid"
    const [uploading, setUploading] = useState(false);
    const [selectedFile, setSelectedFile] = useState(null);
    const fileInputRef = React.useRef(); // Create a ref for the file input
    const navigate = useNavigate(); // Hook to handle navigation

// Load dataset details and filesystem
    const fetchData = async () => {
        try {
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

            // Load filesystem
            const filesystemRes = await axios.get(`${API_URL}/api/v1/nas/pools/naspool/datasets/${datasetId}/files/r`, {
                headers: { Authorization: `${authToken}` },
            });
            setFilesystem(filesystemRes.data.data);
        } catch (error) {
            console.error("Error loading data:", error);
        }
    };

    useEffect(() => {
        fetchData();
    }, [datasetId]);

    const handleDelete = async (path) => {
        // Show confirmation dialog
        const confirmDelete = window.confirm(`Are you sure you want to delete "${path}"?`);

        if (!confirmDelete) {
            return; // Exit if the user cancels
        }

        try {
            // Convert the path to Base64 encoding
            const encodedPath = btoa(path);

            const authToken = localStorage.getItem("auth_token");
            await axios.delete(`${API_URL}/api/v1/nas/pools/naspool/datasets/${datasetId}/files/${encodedPath}`, {
                headers: { Authorization: `${authToken}` },
                data: { path }, // Include the path in the request body
            });
            setFilesystem(filesystem.filter((item) => item.Name !== path));
        } catch (error) {
            console.error("Error deleting file/folder:", error);
        }
    };

    const handleUpload = async () => {
        if (!selectedFile) return;

        const formData = new FormData();
        formData.append("file", selectedFile);

        try {
            const authToken = localStorage.getItem("auth_token");
            setUploading(true);
            await axios.post(`${API_URL}/api/v1/nas/pools/naspool/datasets/${datasetId}/files/r`, formData, {
                headers: {
                    Authorization: `${authToken}`,
                    "Content-Type": "multipart/form-data",
                },
            });
            alert("File uploaded successfully!");
            setSelectedFile(null);

            // Clear the file input after successful upload
            fileInputRef.current.value = ""; // Reset the input field

            fetchData();

        } catch (error) {
            console.error("Error uploading file:", error);
        } finally {
            setUploading(false);
        }
    };

    const renderFilesystemItem = (item) => (
        <div
            key={item.Path}
            style={{
                border: "1px solid #ccc",
                borderRadius: "5px",
                padding: "10px",
                margin: "10px",
                display: "inline-block",
                width: "150px",
                textAlign: "center",
                position: "relative",
            }}
            onContextMenu={(e) => {
                e.preventDefault();
                if (window.confirm(`Are you sure you want to delete ${item.Path}?`)) {
                    handleDelete(item.Path);
                }
            }}
        >
            {item.Size === 0 ? (
                <i className="fas fa-folder" style={{ fontSize: "40px", color: "#f0ad4e" }} />
            ) : (
                <i className="fas fa-file" style={{ fontSize: "40px", color: "#5bc0de" }} />
            )}
            <p style={{ wordWrap: "break-word", fontSize: "14px" }}>{item.Path.split("/").pop()}</p>
        </div>
    );

    return (
        <div style={{ padding: "20px" }}>
            <h2>Filesystem</h2>

            {datasetDetails && (
                <div>
                    <h3>Dataset: {datasetDetails.name}</h3>
                    <p>Quota: {datasetDetails.quota}</p>
                    <p>Used: {datasetDetails.used}</p>
                    <p>Available: {datasetDetails.available}</p>
                </div>
            )}

            {/*<div style={{ marginBottom: "20px" }}>
                <button
                    onClick={() => setViewMode("list")}
                    style={{
                        padding: "10px 20px",
                        backgroundColor: viewMode === "list" ? "#007bff" : "#ccc",
                        color: "#fff",
                        border: "none",
                        borderRadius: "5px",
                        marginRight: "10px",
                    }}
                >
                    List View
                </button>
                <button
                    onClick={() => setViewMode("grid")}
                    style={{
                        padding: "10px 20px",
                        backgroundColor: viewMode === "grid" ? "#007bff" : "#ccc",
                        color: "#fff",
                        border: "none",
                        borderRadius: "5px",
                    }}
                >
                    Grid View
                </button>
            </div>*/}

            <div>
                <input
                    type="file"
                    onChange={(e) => setSelectedFile(e.target.files[0])}
                    style={{ marginBottom: "10px" }}
                />
                <button
                    onClick={handleUpload}
                    ref={fileInputRef} // Link input to the ref
                    disabled={uploading}
                    style={{
                        padding: "10px 20px",
                        backgroundColor: "#28a745",
                        color: "#fff",
                        border: "none",
                        borderRadius: "5px",
                        cursor: uploading ? "not-allowed" : "pointer",
                    }}
                >
                    {uploading ? "Uploading..." : "Upload File"}
                </button>
            </div>

            {filesystem?.length > 0 ? (
                <div style={{ marginTop: "20px" }}>
                    {viewMode === "list" ? (
                        <ul>
                            {filesystem.map((item) => (
                                <li
                                    key={item.Path}
                                    style={{
                                        marginBottom: "10px",
                                        display: "flex",
                                        justifyContent: "space-between",
                                    }}
                                >
                  <span>
                    {item.Size === 0 ? "📁" : "📄"} {item.Name}
                  </span>
                                    <button
                                        onClick={() => handleDelete(item.Name)}
                                        style={{
                                            padding: "5px 10px",
                                            backgroundColor: "#dc3545",
                                            color: "#fff",
                                            border: "none",
                                            borderRadius: "5px",
                                            cursor: "pointer",
                                        }}
                                    >
                                        Delete
                                    </button>
                                </li>
                            ))}
                        </ul>
                    ) : (
                        <div style={{ display: "flex", flexWrap: "wrap" }}>
                            {filesystem.map((item) => renderFilesystemItem(item))}
                        </div>
                    )}
                </div>
            ) : (
                <p>No files or folders found.</p>
            )}
        </div>
    );
};

export default Filesystem;
