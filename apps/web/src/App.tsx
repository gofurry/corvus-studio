import { Layout, Space, Typography } from 'antd';
import { Link, Navigate, Route, Routes } from 'react-router-dom';

import ProjectCreatePage from './features/projects/ProjectCreatePage';
import ProjectDetailPage from './features/projects/ProjectDetailPage';
import ProjectListPage from './features/projects/ProjectListPage';
import ChecklistPage from './features/release/ChecklistPage';
import ReleasePage from './features/release/ReleasePage';
import NotFoundPage from './pages/NotFoundPage';

const { Header, Content } = Layout;

function App() {
  return (
    <Layout className="app-shell">
      <Header className="app-header">
        <Link className="app-brand" to="/projects" aria-label="Corvus Studio projects">
          <span className="app-brand-mark">C</span>
          <Space orientation="vertical" size={0}>
            <Typography.Text className="app-brand-name">Corvus Studio</Typography.Text>
            <Typography.Text className="app-brand-caption">Local release workspace</Typography.Text>
          </Space>
        </Link>
      </Header>
      <Content className="app-content">
        <Routes>
          <Route path="/" element={<Navigate to="/projects" replace />} />
          <Route path="/projects" element={<ProjectListPage />} />
          <Route path="/projects/new" element={<ProjectCreatePage />} />
          <Route path="/projects/:projectId" element={<ProjectDetailPage />} />
          <Route path="/projects/:projectId/release" element={<ReleasePage />} />
          <Route path="/projects/:projectId/checklist" element={<ChecklistPage />} />
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </Content>
    </Layout>
  );
}

export default App;
