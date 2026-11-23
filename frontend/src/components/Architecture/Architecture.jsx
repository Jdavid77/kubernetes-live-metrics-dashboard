import React, { useState, useEffect } from 'react';
import yaml from 'js-yaml';
import FlowDiagram from './FlowDiagram';
import LoadingSpinner from '../Common/LoadingSpinner';
import './Architecture.css';

// Loaded at runtime so the content can be changed without rebuilding the image:
// mount a ConfigMap over /usr/share/nginx/html/config/ to override the default.
const TOPICS_URL = `${process.env.PUBLIC_URL}/config/architecture.yaml`;

const Architecture = () => {
  const [topics, setTopics] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchTopics = async () => {
      try {
        const res = await fetch(TOPICS_URL, { cache: 'no-cache' });
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        const data = yaml.load(await res.text());
        // nginx falls back to index.html for missing files, so check the shape.
        if (!Array.isArray(data?.topics)) throw new Error('missing "topics" list');
        setTopics(data.topics);
        setError(null);
      } catch (err) {
        console.error('Error loading architecture topics:', err);
        setError('Failed to load architecture');
      } finally {
        setLoading(false);
      }
    };

    fetchTopics();
  }, []);

  if (loading) {
    return <LoadingSpinner />;
  }

  if (error) {
    return <div className="text-center text-muted">{error}</div>;
  }

  return (
    <div className="architecture">
      <div className="container">
        <nav className="architecture-toc">
          {topics.map((topic) => (
            <a key={topic.id} href={`#${topic.id}`} style={{ '--accent': topic.color }}>
              {topic.title}
            </a>
          ))}
        </nav>

        {topics.map((topic) => (
          <section
            key={topic.id}
            id={topic.id}
            className="topic-card"
            style={{ '--accent': topic.color }}
          >
            <h2 className="topic-title">{topic.title}</h2>
            <p className="topic-summary">{topic.summary}</p>

            {topic.lanes && <FlowDiagram lanes={topic.lanes} />}

            {topic.components && (
              <>
                <h3 className="topic-subtitle">Components</h3>
                <ul className="topic-components">
                  {topic.components.map((c) => (
                    <li key={c.name}>
                      <strong>{c.name}</strong>
                      <span>{c.role}</span>
                    </li>
                  ))}
                </ul>
              </>
            )}

            {topic.decisions && (
              <>
                <h3 className="topic-subtitle">Decisions</h3>
                <div className="topic-decisions">
                  {topic.decisions.map((d) => (
                    <div key={d.title} className="decision">
                      <h4>{d.title}</h4>
                      <p>{d.why}</p>
                    </div>
                  ))}
                </div>
              </>
            )}
          </section>
        ))}
      </div>
    </div>
  );
};

export default Architecture;
