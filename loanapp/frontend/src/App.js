import React, { useState, useEffect } from 'react';
import './App.css';

const API_URL = 'http://localhost:8080/api';

function App() {
  const [applications, setApplications] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [filter, setFilter] = useState('all'); // all, approved, review_required, rejected

  useEffect(() => {
    fetchApplications();
  }, []);

  // Set up long polling
  useEffect(() => {
    let isActive = true;

    const startLongPolling = async () => {
      while (isActive) {
        try {
          const response = await fetch(`${API_URL}/poll`);
          if (!response.ok) {
            console.error('Long polling error:', response.statusText);
            // Wait a bit before retrying on error
            await new Promise(resolve => setTimeout(resolve, 5000));
            continue;
          }

          const data = await response.json();

          if (data.action === 'refresh' && isActive) {
            // Server told us to refresh
            console.log('Received refresh signal from server');
            fetchApplications();
          }

          // Continue polling (either after refresh or timeout)
        } catch (err) {
          console.error('Long polling error:', err);
          // Wait a bit before retrying on error
          await new Promise(resolve => setTimeout(resolve, 5000));
        }
      }
    };

    startLongPolling();

    // Cleanup function to stop polling when component unmounts
    return () => {
      isActive = false;
    };
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  const fetchApplications = async () => {
    try {
      setLoading(true);
      const response = await fetch(`${API_URL}/applications`);
      if (!response.ok) {
        throw new Error('Failed to fetch applications');
      }
      const data = await response.json();
      setApplications(data.applications || []);
      setError(null);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const getRiskLevel = (riskFactor) => {
    if (riskFactor < 0.3) return 'low';
    if (riskFactor < 0.7) return 'medium';
    return 'high';
  };

  const getStatusClass = (status) => {
    return `status status-${status}`;
  };

  const getRiskClass = (riskFactor) => {
    return `risk risk-${getRiskLevel(riskFactor)}`;
  };

  const filteredApplications = applications.filter(app => {
    if (filter === 'all') return true;
    return app.status === filter;
  });

  if (loading) {
    return (
      <div className="App">
        <div className="loading">Loading applications...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="App">
        <div className="error">
          <h2>Error</h2>
          <p>{error}</p>
          <button onClick={fetchApplications}>Retry</button>
        </div>
      </div>
    );
  }

  return (
    <div className="App">
      <header className="App-header">
        <h1>LoanApp - Loan Applications</h1>
        <p className="subtitle">Reviewing {applications.length} loan applications</p>
      </header>

      <div className="filters">
        <button
          className={filter === 'all' ? 'active' : ''}
          onClick={() => setFilter('all')}
        >
          All ({applications.length})
        </button>
        <button
          className={filter === 'approved' ? 'active' : ''}
          onClick={() => setFilter('approved')}
        >
          Approved ({applications.filter(a => a.status === 'approved').length})
        </button>
        <button
          className={filter === 'review_required' ? 'active' : ''}
          onClick={() => setFilter('review_required')}
        >
          Review Required ({applications.filter(a => a.status === 'review_required').length})
        </button>
        <button
          className={filter === 'rejected' ? 'active' : ''}
          onClick={() => setFilter('rejected')}
        >
          Rejected ({applications.filter(a => a.status === 'rejected').length})
        </button>
      </div>

      <div className="applications-container">
        {filteredApplications.length === 0 ? (
          <div className="no-results">No applications found</div>
        ) : (
          <table className="applications-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Applicant Name</th>
                <th>Credit Score</th>
                <th>Amount Requested</th>
                <th>Risk Factor</th>
                <th>Risk Level</th>
                <th>Status</th>
                <th>Created At</th>
              </tr>
            </thead>
            <tbody>
              {filteredApplications.map((app) => (
                <tr key={app.id}>
                  <td>{app.id}</td>
                  <td className="name">{app.applicant_name}</td>
                  <td className="credit-score">{app.credit_score}</td>
                  <td className="amount">${app.amount_requested.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</td>
                  <td className="risk-factor">{(app.risk_factor * 100).toFixed(1)}%</td>
                  <td>
                    <span className={getRiskClass(app.risk_factor)}>
                      {getRiskLevel(app.risk_factor).toUpperCase()}
                    </span>
                  </td>
                  <td>
                    <span className={getStatusClass(app.status)}>
                      {app.status.toUpperCase()}
                    </span>
                  </td>
                  <td className="date">
                    {new Date(app.created_at).toLocaleDateString('en-US', {
                      month: 'short',
                      day: 'numeric',
                      year: 'numeric'
                    })}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}

export default App;
