import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  Alert,
  Button,
  Card,
  Flex,
  Form,
  Input,
  InputNumber,
  Select,
  Space,
  Typography,
} from 'antd';
import { Link, useNavigate } from 'react-router-dom';

import { projectsApi, type CreateProjectRequest, type ProjectStage } from '../../api/projects';
import { projectStageLabels } from './projectPresentation';
import { projectQueryKeys } from './projectQueries';

const stageOptions = (Object.keys(projectStageLabels) as ProjectStage[]).map((value) => ({
  value,
  label: projectStageLabels[value],
}));

function ProjectCreatePage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const createProject = useMutation({
    mutationFn: projectsApi.create,
    onSuccess: async (project) => {
      await queryClient.invalidateQueries({ queryKey: projectQueryKeys.all });
      queryClient.setQueryData(projectQueryKeys.detail(project.id), project);
      navigate(`/projects/${project.id}`);
    },
  });

  return (
    <section className="page-stack narrow-page" aria-labelledby="create-project-heading">
      <div>
        <Typography.Text className="page-eyebrow">New workspace</Typography.Text>
        <Typography.Title id="create-project-heading" level={1}>
          Create project
        </Typography.Title>
        <Typography.Paragraph className="page-intro">
          Corvus Studio stores a reference to an existing directory. It will not create, move, or
          modify files there.
        </Typography.Paragraph>
      </div>

      <Card className="form-card">
        {createProject.isError ? (
          <Alert
            className="form-alert"
            type="error"
            showIcon
            title="Project could not be created"
            description={createProject.error.message}
          />
        ) : null}
        <Form<CreateProjectRequest>
          layout="vertical"
          initialValues={{ language: 'English', stage: 'concept' }}
          onFinish={(values) => createProject.mutate(values)}
          requiredMark="optional"
        >
          <Form.Item
            label="Game name"
            name="name"
            rules={[
              { required: true, whitespace: true, message: 'Enter a game name' },
              { max: 120 },
            ]}
          >
            <Input autoFocus maxLength={120} placeholder="Raven Game" />
          </Form.Item>

          <Form.Item label="Description" name="description" rules={[{ max: 2000 }]}>
            <Input.TextArea
              maxLength={2000}
              placeholder="A short note about this project"
              rows={3}
              showCount
            />
          </Form.Item>

          <Form.Item
            label="Project location"
            name="location"
            extra="Enter an existing absolute directory on this computer."
            rules={[{ required: true, whitespace: true, message: 'Enter an absolute directory' }]}
          >
            <Input placeholder="C:\Games\Raven or /home/me/games/raven" />
          </Form.Item>

          <Flex gap="middle" wrap>
            <Form.Item
              className="form-flex-item"
              label="Primary language"
              name="language"
              rules={[
                { required: true, whitespace: true, message: 'Enter the primary language' },
                { max: 64 },
              ]}
            >
              <Input maxLength={64} />
            </Form.Item>
            <Form.Item
              className="form-flex-item"
              label="Current stage"
              name="stage"
              rules={[{ required: true }]}
            >
              <Select options={stageOptions} />
            </Form.Item>
          </Flex>

          <Form.Item
            label="Steam AppID"
            name="steam_app_id"
            extra="Optional. Use the positive numeric AppID assigned by Steam."
          >
            <InputNumber min={1} max={4294967295} precision={0} style={{ width: '100%' }} />
          </Form.Item>

          <Space>
            <Button type="primary" htmlType="submit" loading={createProject.isPending}>
              Create project
            </Button>
            <Link to="/projects">
              <Button>Cancel</Button>
            </Link>
          </Space>
        </Form>
      </Card>
    </section>
  );
}

export default ProjectCreatePage;
