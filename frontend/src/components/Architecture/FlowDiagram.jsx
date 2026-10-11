import React from 'react';
import { FiArrowRight } from 'react-icons/fi';

const FlowDiagram = ({ lanes }) => (
  <div className="flow-diagram">
    {lanes.map((lane) => (
      <div key={lane.label} className="flow-lane">
        <span className="flow-lane-label">{lane.label}</span>
        <div className="flow-steps">
          {lane.steps.map((step, i) => (
            <React.Fragment key={step}>
              {i > 0 && <FiArrowRight className="flow-arrow" aria-hidden="true" />}
              <span className="flow-step">{step}</span>
            </React.Fragment>
          ))}
        </div>
      </div>
    ))}
  </div>
);

export default FlowDiagram;
