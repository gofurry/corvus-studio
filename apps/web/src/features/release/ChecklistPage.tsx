import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  Alert,
  Button,
  Card,
  Drawer,
  Empty,
  Flex,
  Form,
  Grid,
  Input,
  List,
  Modal,
  Select,
  Skeleton,
  Space,
  Tag,
  Typography,
} from 'antd';
import { useState } from 'react';
import { Link, useParams, useSearchParams } from 'react-router-dom';

import {
  checklistApi,
  type ChecklistCategory,
  type ChecklistFilters,
  type ChecklistItem,
  type ChecklistRequirementLevel,
  type ChecklistSource,
  type ChecklistStatus,
  type CreateChecklistItemRequest,
} from '../../api/checklist';
import { projectsApi } from '../../api/projects';
import { releasesApi } from '../../api/releases';
import ProjectWorkspaceNav from '../projects/ProjectWorkspaceNav';
import { projectQueryKeys } from '../projects/projectQueries';
import {
  categoryLabels,
  checklistStatusLabels,
  releaseStatusLabels,
  requirementLevelLabels,
  sourceLabels,
} from './workflowPresentation';
import { checklistQueryKeys, releaseQueryKeys } from './workflowQueries';

const checklistStatuses = Object.keys(checklistStatusLabels) as ChecklistStatus[];
const checklistSources: ChecklistSource[] = ['platform_template', 'corvus_template', 'user'];
const checklistCategories = Object.keys(categoryLabels) as ChecklistCategory[];

function ChecklistPage() {
  const { projectId = '' } = useParams();
  const [searchParams, setSearchParams] = useSearchParams();
  const [customTaskOpen, setCustomTaskOpen] = useState(false);
  const [form] = Form.useForm<CustomTaskForm>();
  const screens = Grid.useBreakpoint();
  const desktop = screens.md === true;
  const queryClient = useQueryClient();

  const project = useQuery({
    queryKey: projectQueryKeys.detail(projectId),
    queryFn: () => projectsApi.get(projectId),
    enabled: projectId !== '',
  });
  const releases = useQuery({
    queryKey: releaseQueryKeys.project(projectId),
    queryFn: () => releasesApi.listForProject(projectId),
    enabled: projectId !== '',
  });
  const release = releases.data?.find((candidate) => candidate.goal_type === 'steam_coming_soon');
  const filters = filtersFromSearch(searchParams);
  const selectedItemId = searchParams.get('item') ?? '';

  const checklist = useQuery({
    queryKey: checklistQueryKeys.list(release?.id ?? '', filters),
    queryFn: () => checklistApi.list(release?.id ?? '', filters),
    enabled: release !== undefined,
  });
  const selectedItem = useQuery({
    queryKey: checklistQueryKeys.detail(selectedItemId),
    queryFn: () => checklistApi.get(selectedItemId),
    enabled: selectedItemId !== '',
  });

  const transition = useMutation({
    mutationFn: ({ itemId, status }: { itemId: string; status: ChecklistStatus }) =>
      checklistApi.transition(itemId, status),
    onSuccess: (updated) => {
      queryClient.setQueryData(checklistQueryKeys.detail(updated.id), updated);
      if (release !== undefined) {
        void queryClient.invalidateQueries({ queryKey: ['checklist', release.id] });
        void queryClient.invalidateQueries({ queryKey: releaseQueryKeys.project(projectId) });
        void queryClient.invalidateQueries({ queryKey: releaseQueryKeys.detail(release.id) });
      }
    },
  });

  const createTask = useMutation({
    mutationFn: (values: CustomTaskForm) => {
      if (release === undefined) {
        throw new Error('Create a Release Goal before adding tasks');
      }
      const input: CreateChecklistItemRequest = {
        release_goal_id: release.id,
        title: values.title,
        description: values.description ?? '',
        requirement: values.requirement ?? '',
        category: values.category,
        requirement_level: values.requirement_level,
      };
      return checklistApi.create(input);
    },
    onSuccess: (created) => {
      if (release !== undefined) {
        void queryClient.invalidateQueries({ queryKey: ['checklist', release.id] });
        void queryClient.invalidateQueries({ queryKey: releaseQueryKeys.project(projectId) });
      }
      form.resetFields();
      setCustomTaskOpen(false);
      updateSearchParam(searchParams, setSearchParams, 'item', created.id);
    },
  });

  if (project.isPending || releases.isPending) {
    return (
      <section className="page-stack" aria-label="Loading checklist workspace">
        <Card className="detail-card">
          <Skeleton active />
        </Card>
      </section>
    );
  }

  const loadError = project.error ?? releases.error;
  if (loadError !== null) {
    return (
      <section className="page-stack">
        <Alert
          type="error"
          showIcon
          title="Checklist workspace is unavailable"
          description={loadError.message}
        />
      </section>
    );
  }

  if (project.data === undefined) {
    return null;
  }

  return (
    <section className="page-stack" aria-labelledby="checklist-heading">
      <Flex justify="space-between" align="flex-start" gap="middle" wrap>
        <div>
          <Typography.Text className="page-eyebrow">{project.data.name}</Typography.Text>
          <Typography.Title id="checklist-heading" level={1}>
            Release checklist
          </Typography.Title>
          <Typography.Paragraph className="page-intro">
            Filter the work, inspect each requirement, and move tasks forward without losing
            context.
          </Typography.Paragraph>
        </div>
        <Link to="/projects">
          <Button>All projects</Button>
        </Link>
      </Flex>

      <ProjectWorkspaceNav projectId={projectId} />

      {release === undefined ? (
        <Card className="empty-card center-state">
          <Empty description="Create the Steam Coming Soon Release Goal before opening its Checklist.">
            <Link to={`/projects/${projectId}/release`}>
              <Button type="primary">Set up release workspace</Button>
            </Link>
          </Empty>
        </Card>
      ) : (
        <>
          <Flex justify="space-between" align="center" gap="middle" wrap>
            <Space wrap>
              <Tag color="geekblue">{releaseStatusLabels[release.status]}</Tag>
              <Typography.Text type="secondary">
                {release.checklist_summary.done}/{release.checklist_summary.total} completed
              </Typography.Text>
            </Space>
            <Button type="primary" onClick={() => setCustomTaskOpen(true)}>
              Add custom task
            </Button>
          </Flex>

          {transition.error !== null && (
            <Alert
              type="warning"
              showIcon
              title="Task status was not changed"
              description={transition.error.message}
            />
          )}
          {checklist.isError && (
            <Alert
              type="error"
              showIcon
              title="Checklist could not be loaded"
              description={checklist.error.message}
            />
          )}

          <Card className="detail-card checklist-filters" title="Filters">
            <Space wrap>
              <Select
                aria-label="Filter by status"
                allowClear
                placeholder="All statuses"
                value={filters.status}
                options={checklistStatuses.map((value) => ({
                  value,
                  label: checklistStatusLabels[value],
                }))}
                onChange={(value) =>
                  updateSearchParam(searchParams, setSearchParams, 'status', value)
                }
              />
              <Select
                aria-label="Filter by source"
                allowClear
                placeholder="All sources"
                value={filters.source}
                options={checklistSources.map((value) => ({ value, label: sourceLabels[value] }))}
                onChange={(value) =>
                  updateSearchParam(searchParams, setSearchParams, 'source', value)
                }
              />
              <Select
                aria-label="Filter by category"
                allowClear
                placeholder="All categories"
                value={filters.category}
                options={checklistCategories.map((value) => ({
                  value,
                  label: categoryLabels[value],
                }))}
                onChange={(value) =>
                  updateSearchParam(searchParams, setSearchParams, 'category', value)
                }
              />
            </Space>
          </Card>

          <div className="checklist-workspace">
            <Card className="detail-card checklist-categories" title="Categories">
              <Space orientation="vertical" className="full-width">
                <Button
                  type={filters.category === undefined ? 'primary' : 'text'}
                  block
                  onClick={() => updateSearchParam(searchParams, setSearchParams, 'category')}
                >
                  All categories
                </Button>
                {checklistCategories.map((category) => (
                  <Button
                    key={category}
                    type={filters.category === category ? 'primary' : 'text'}
                    block
                    onClick={() =>
                      updateSearchParam(searchParams, setSearchParams, 'category', category)
                    }
                  >
                    {categoryLabels[category]}
                  </Button>
                ))}
              </Space>
            </Card>

            <Card className="detail-card checklist-list-card" title="Tasks">
              <List
                loading={checklist.isPending}
                locale={{ emptyText: 'No tasks match these filters.' }}
                dataSource={checklist.data ?? []}
                renderItem={(item) => (
                  <List.Item>
                    <button
                      className={`checklist-task${selectedItemId === item.id ? ' is-selected' : ''}`}
                      type="button"
                      onClick={() =>
                        updateSearchParam(searchParams, setSearchParams, 'item', item.id)
                      }
                    >
                      <span>{item.title}</span>
                      <Space wrap size="small">
                        <Tag>{categoryLabels[item.category]}</Tag>
                        <Tag color={item.status === 'done' ? 'green' : 'default'}>
                          {checklistStatusLabels[item.status]}
                        </Tag>
                      </Space>
                    </button>
                  </List.Item>
                )}
              />
            </Card>

            {desktop && (
              <Card className="detail-card checklist-detail-card" title="Task details">
                <TaskDetail
                  item={selectedItem.data}
                  loading={selectedItem.isPending && selectedItemId !== ''}
                  error={selectedItem.error}
                  transitioning={transition.isPending}
                  onTransition={(status) => {
                    if (selectedItem.data !== undefined) {
                      transition.mutate({ itemId: selectedItem.data.id, status });
                    }
                  }}
                />
              </Card>
            )}
          </div>

          {!desktop && (
            <Drawer
              title="Task details"
              open={selectedItemId !== ''}
              onClose={() => updateSearchParam(searchParams, setSearchParams, 'item')}
              size="large"
            >
              <TaskDetail
                item={selectedItem.data}
                loading={selectedItem.isPending}
                error={selectedItem.error}
                transitioning={transition.isPending}
                onTransition={(status) => {
                  if (selectedItem.data !== undefined) {
                    transition.mutate({ itemId: selectedItem.data.id, status });
                  }
                }}
              />
            </Drawer>
          )}

          <CustomTaskModal
            open={customTaskOpen}
            form={form}
            saving={createTask.isPending}
            error={createTask.error}
            onCancel={() => setCustomTaskOpen(false)}
            onFinish={(values) => createTask.mutate(values)}
          />
        </>
      )}
    </section>
  );
}

type TaskDetailProps = {
  item?: ChecklistItem;
  loading: boolean;
  error: Error | null;
  transitioning: boolean;
  onTransition: (status: ChecklistStatus) => void;
};

function TaskDetail({ item, loading, error, transitioning, onTransition }: TaskDetailProps) {
  if (loading) {
    return <Skeleton active />;
  }
  if (error !== null) {
    return (
      <Alert type="error" showIcon title="Task could not be opened" description={error.message} />
    );
  }
  if (item === undefined) {
    return <Empty description="Select a task to inspect its details." />;
  }

  const statuses =
    item.requirement_level === 'required'
      ? checklistStatuses.filter((status) => status !== 'not_applicable')
      : checklistStatuses;

  return (
    <Space orientation="vertical" size="large" className="full-width">
      <div>
        <Typography.Title level={3}>{item.title}</Typography.Title>
        <Space wrap>
          <Tag color={item.requirement_level === 'required' ? 'gold' : 'default'}>
            {requirementLevelLabels[item.requirement_level]}
          </Tag>
          <Tag>{sourceLabels[item.source]}</Tag>
        </Space>
      </div>
      <section>
        <Typography.Text strong>What</Typography.Text>
        <Typography.Paragraph>
          {item.description || 'No description provided.'}
        </Typography.Paragraph>
      </section>
      <section>
        <Typography.Text strong>Why</Typography.Text>
        <Typography.Paragraph>
          {item.requirement || 'No requirement note provided.'}
        </Typography.Paragraph>
      </section>
      <section>
        <Typography.Text strong>Requirement</Typography.Text>
        <Typography.Paragraph>
          {requirementLevelLabels[item.requirement_level]} · {categoryLabels[item.category]}
        </Typography.Paragraph>
      </section>
      <section>
        <Typography.Text strong>Source</Typography.Text>
        <Typography.Paragraph>
          {item.source_reference.startsWith('http') ? (
            <a href={item.source_reference} target="_blank" rel="noreferrer">
              {sourceLabels[item.source]} reference
            </a>
          ) : (
            item.source_reference || sourceLabels[item.source]
          )}
        </Typography.Paragraph>
      </section>
      <section>
        <Typography.Text strong>Status</Typography.Text>
        <Space wrap className="status-actions">
          {statuses.map((status) => (
            <Button
              key={status}
              type={status === item.status ? 'primary' : 'default'}
              disabled={status === item.status}
              loading={transitioning}
              aria-label={`Set task status to ${checklistStatusLabels[status]}`}
              onClick={() => onTransition(status)}
            >
              {checklistStatusLabels[status]}
            </Button>
          ))}
        </Space>
      </section>
    </Space>
  );
}

type CustomTaskForm = {
  title: string;
  description?: string;
  requirement?: string;
  category: ChecklistCategory;
  requirement_level: ChecklistRequirementLevel;
};

type CustomTaskModalProps = {
  open: boolean;
  form: ReturnType<typeof Form.useForm<CustomTaskForm>>[0];
  saving: boolean;
  error: Error | null;
  onCancel: () => void;
  onFinish: (values: CustomTaskForm) => void;
};

function CustomTaskModal({ open, form, saving, error, onCancel, onFinish }: CustomTaskModalProps) {
  return (
    <Modal
      title="Add custom task"
      open={open}
      confirmLoading={saving}
      okText="Add task"
      onCancel={onCancel}
      onOk={() => form.submit()}
      destroyOnHidden
    >
      {error !== null && (
        <Alert
          className="form-alert"
          type="error"
          showIcon
          title="Custom task was not added"
          description={error.message}
        />
      )}
      <Form<CustomTaskForm>
        form={form}
        layout="vertical"
        initialValues={{ category: 'review', requirement_level: 'recommended' }}
        onFinish={onFinish}
      >
        <Form.Item
          name="title"
          label="Task title"
          rules={[{ required: true, whitespace: true, max: 160 }]}
        >
          <Input maxLength={160} />
        </Form.Item>
        <Form.Item name="description" label="What" rules={[{ max: 2000 }]}>
          <Input.TextArea rows={3} maxLength={2000} />
        </Form.Item>
        <Form.Item name="requirement" label="Why" rules={[{ max: 2000 }]}>
          <Input.TextArea rows={3} maxLength={2000} />
        </Form.Item>
        <Form.Item name="category" label="Category" rules={[{ required: true }]}>
          <Select
            options={checklistCategories.map((category) => ({
              value: category,
              label: categoryLabels[category],
            }))}
          />
        </Form.Item>
        <Form.Item name="requirement_level" label="Requirement level" rules={[{ required: true }]}>
          <Select
            options={(['required', 'recommended'] as ChecklistRequirementLevel[]).map((level) => ({
              value: level,
              label: requirementLevelLabels[level],
            }))}
          />
        </Form.Item>
      </Form>
    </Modal>
  );
}

function filtersFromSearch(searchParams: URLSearchParams): ChecklistFilters {
  const status = searchParams.get('status');
  const source = searchParams.get('source');
  const category = searchParams.get('category');
  return {
    status: checklistStatuses.includes(status as ChecklistStatus)
      ? (status as ChecklistStatus)
      : undefined,
    source: checklistSources.includes(source as ChecklistSource)
      ? (source as ChecklistSource)
      : undefined,
    category: checklistCategories.includes(category as ChecklistCategory)
      ? (category as ChecklistCategory)
      : undefined,
  };
}

function updateSearchParam(
  current: URLSearchParams,
  setSearchParams: ReturnType<typeof useSearchParams>[1],
  key: string,
  value?: string,
) {
  const next = new URLSearchParams(current);
  if (value === undefined || value === '') {
    next.delete(key);
  } else {
    next.set(key, value);
  }
  setSearchParams(next);
}

export default ChecklistPage;
