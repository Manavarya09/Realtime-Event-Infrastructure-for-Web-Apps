import React from 'react';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import Container from '@mui/material/Container';
import Typography from '@mui/material/Typography';
import Grid from '@mui/material/Grid';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

const theme = createTheme();

const mockData = [
  { time: '00:00', events: 120 },
  { time: '01:00', events: 150 },
  { time: '02:00', events: 180 },
  // ... more data
];

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Realtime Event Analytics Dashboard
        </Typography>

        <Grid container spacing={3}>
          {/* Event Count Card */}
          <Grid item xs={12} md={4}>
            <Card>
              <CardContent>
                <Typography color="textSecondary" gutterBottom>
                  Total Events (Last 24h)
                </Typography>
                <Typography variant="h5">
                  12,456
                </Typography>
              </CardContent>
            </Card>
          </Grid>

          {/* Unique Users Card */}
          <Grid item xs={12} md={4}>
            <Card>
              <CardContent>
                <Typography color="textSecondary" gutterBottom>
                  Unique Users
                </Typography>
                <Typography variant="h5">
                  8,921
                </Typography>
              </CardContent>
            </Card>
          </Grid>

          {/* Active Webhooks Card */}
          <Grid item xs={12} md={4}>
            <Card>
              <CardContent>
                <Typography color="textSecondary" gutterBottom>
                  Active Webhooks
                </Typography>
                <Typography variant="h5">
                  23
                </Typography>
              </CardContent>
            </Card>
          </Grid>

          {/* Events Chart */}
          <Grid item xs={12}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  Events Over Time
                </Typography>
                <ResponsiveContainer width="100%" height={300}>
                  <LineChart data={mockData}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="time" />
                    <YAxis />
                    <Tooltip />
                    <Line type="monotone" dataKey="events" stroke="#8884d8" />
                  </LineChart>
                </ResponsiveContainer>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      </Container>
    </ThemeProvider>
  );
}

export default App;