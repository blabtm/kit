import React from 'react';
import { withJsonFormsLayoutProps } from '@jsonforms/react';
import { rankWith } from '@jsonforms/core';

const Dashboard = ({ uischema }) => {
  const url = uischema.url; 

  if (!url) return null;

  return (
    <div className="lab-ui-iframe" style={{ margin: '4px 0', lineHeight: 0 }}>
      <iframe
        src={`http://192.168.0.37:3000/${url}`}
        width="100%"
        height="300"
        title="Dashboard"
      />
    </div>
  );
};

export const dashboardRenderer = withJsonFormsLayoutProps(Dashboard);
export const dashboardTester = rankWith(
  10,
  uischema => uischema.type === 'Dashboard'
);
