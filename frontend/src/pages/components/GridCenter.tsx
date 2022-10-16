import { Grid } from '@mui/material';
import React from 'react';

export const GridCenter = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  return (
    <Grid
      container
      marginTop={1}
      display="flex"
      justifyContent="center"
      alignItems="center"
    >
      {children}
    </Grid>
  );
};
