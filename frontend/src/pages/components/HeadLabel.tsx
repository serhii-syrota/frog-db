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
export const HeadLabelClickable = (props: {
  children: React.ReactNode;
  href: string;
}) => {
  return (
    <Typography
      href={props.href}
      style={{ paddingLeft: 10 }}
      variant="h6"
      component="a"
      noWrap
      color="inherit"
      sx={{
        textDecoration: 'none',
      }}
    >
      {props.children}
    </Typography>
  );
};
