import { Typography } from '@mui/material';
import React from 'react';

export const HeadLabel = (props: { children: React.ReactNode }) => {
  return (
    <Typography
      style={{ paddingLeft: 10 }}
      variant="h6"
      color="inherit"
      noWrap
    >
      {props.children}
    </Typography>
  );
};
