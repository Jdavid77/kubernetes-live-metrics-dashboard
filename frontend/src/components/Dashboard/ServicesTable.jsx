import React, { useState, useEffect } from 'react';
import { formatRelativeTime } from '../../utils/formatters';
import { getServices } from '../../services/servicesService';
import LoadingSpinner from '../Common/LoadingSpinner';
import './ResourceTable.css';

const ServicesTable = ({ namespace }) => {
  const [services, setServices] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchServices = async () => {
      try {
        setLoading(true);
        const data = await getServices(namespace);
        setServices(data || []);
        setError(null);
      } catch (err) {
        console.error('Error fetching services:', err);
        setError('Failed to fetch services');
      } finally {
        setLoading(false);
      }
    };

    fetchServices();
  }, [namespace]);

  if (loading) {
    return <LoadingSpinner />;
  }

  if (error) {
    return <div className="text-center text-muted">{error}</div>;
  }

  return (
    <div className="resource-table-container">
      <div className="table-header">
        <h2 className="section-title">Services</h2>
        <div className="table-stats">
          Total: {services.length} services
        </div>
      </div>

      <div className="table-wrapper">
        <table className="resource-table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Namespace</th>
              <th>Type</th>
              <th>Cluster IP</th>
              <th>External IP</th>
              <th>Ports</th>
              <th>Age</th>
            </tr>
          </thead>
          <tbody>
            {services.map((service, index) => (
              <tr key={index}>
                <td className="deployment-name">{service.name}</td>
                <td>
                  <span className="namespace-badge">{service.namespace}</span>
                </td>
                <td>
                  <span className="service-type-badge">{service.type}</span>
                </td>
                <td className="ip-cell">{service.cluster_ip}</td>
                <td className="ip-cell">
                  {service.external_ip || <span className="text-muted">None</span>}
                </td>
                <td className="ports-cell">
                  {service.ports && service.ports.length > 0 ? (
                    service.ports.map((port, i) => (
                      <div key={i} className="port-item">
                        {port.port}/{port.protocol}
                      </div>
                    ))
                  ) : (
                    <span className="text-muted">-</span>
                  )}
                </td>
                <td className="age-cell">
                  {formatRelativeTime(service.created_at)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default ServicesTable;
