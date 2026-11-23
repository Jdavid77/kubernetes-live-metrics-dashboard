import React, { useState } from 'react';
import { FiChevronRight } from 'react-icons/fi';

// One swimlane per lane. Hovering a step highlights the same step in every lane.
const FlowDiagram = ({ lanes }) => {
  const [hovered, setHovered] = useState(null);

  return (
    <div className={`flow-diagram${hovered ? ' has-hover' : ''}`}>
      {lanes.map((lane, laneIndex) => (
        <div key={lane.label} className="flow-lane" style={{ '--lane-index': laneIndex }}>
          <span className="flow-lane-label">{lane.label}</span>
          <div className="flow-track">
            <span className="flow-pulse" aria-hidden="true" />
            {lane.steps.map((step, i) => (
              <React.Fragment key={step}>
                {i > 0 && <FiChevronRight className="flow-arrow" aria-hidden="true" />}
                <span
                  className={`flow-step${hovered === step ? ' is-highlighted' : ''}`}
                  onMouseEnter={() => setHovered(step)}
                  onMouseLeave={() => setHovered(null)}
                >
                  {step}
                </span>
              </React.Fragment>
            ))}
          </div>
        </div>
      ))}
    </div>
  );
};

export default FlowDiagram;
