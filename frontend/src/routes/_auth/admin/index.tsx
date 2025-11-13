import { createFileRoute } from '@tanstack/react-router';

import {
  Bar,
  BarChart,
  CartesianGrid,
  LabelList,
  Pie,
  PieChart,
  XAxis,
  YAxis,
} from 'recharts';
import type { PieLabelRenderProps } from 'recharts/types/polar/Pie';

import {
  type MaterialType,
  type TopBusiness,
  getApiAnalyticsMaterialTypes,
  getApiAnalyticsTopBusinesses,
  getApiAnalyticsTotalCollectors,
  getApiAnalyticsTotalWasteCollected,
} from '@/api-client';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
} from '@/components/ui/card';
import {
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
} from '@/components/ui/chart';
import { apiClient } from '@/lib/utils';

const renderCustomizedLabel = ({
  cx,
  cy,
  midAngle,
  innerRadius,
  outerRadius,
  name,
}: PieLabelRenderProps) => {
  const RADIAN = Math.PI / 180;
  const radius =
    Number(innerRadius) + (Number(outerRadius) - Number(innerRadius)) * 1.2;
  const x = Number(cx) + Number(radius) * Math.cos(-Number(midAngle) * RADIAN);
  const y = Number(cy) + Number(radius) * Math.sin(-Number(midAngle) * RADIAN);

  // You can set the color dynamically based on 'index' or other data properties
  const labelColor = 'var(--foreground)';

  return (
    <text
      x={x}
      y={y}
      fill={labelColor}
      textAnchor={Number(x) > Number(cx) ? 'start' : 'end'}
      dominantBaseline="central"
    >
      {name}
    </text>
  );
};

export const Route = createFileRoute('/_auth/admin/')({
  loader: async () => {
    const { data: totalWasteCollectedData } =
      await getApiAnalyticsTotalWasteCollected({
        client: apiClient,
        throwOnError: true,
      });

    const { data: totalCollectorsData } = await getApiAnalyticsTotalCollectors({
      client: apiClient,
      throwOnError: true,
    });

    const { data: topBusinessesData } = await getApiAnalyticsTopBusinesses({
      client: apiClient,
      throwOnError: true,
    });

    const { data: materialTypesData } = await getApiAnalyticsMaterialTypes({
      client: apiClient,
      throwOnError: true,
    });

    return {
      totalWasteCollected: (totalWasteCollectedData.item ?? 0) as number,
      totalCollectors: (totalCollectorsData.item ?? 0) as number,
      topBusinesses: (topBusinessesData.items ?? []) as TopBusiness[],
      materialTypes: (materialTypesData.items ?? []) as MaterialType[],
    };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const { totalWasteCollected, totalCollectors, topBusinesses, materialTypes } =
    Route.useLoaderData();

  return (
    <div className="flex flex-col w-full h-full gap-3">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
        <Card className="bg-background">
          <CardHeader>
            <CardDescription>Total Waste Collected (kg)</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col items-center justify-center">
              <span className="text-6xl font-bold font-mono">
                {totalWasteCollected}
              </span>
            </div>
          </CardContent>
        </Card>

        <Card className="bg-background">
          <CardHeader>
            <CardDescription>Number of Collectors</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col items-center justify-center">
              <span className="text-6xl font-bold font-mono">
                {totalCollectors}
              </span>
            </div>
          </CardContent>
        </Card>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
        <Card className="bg-background">
          <CardHeader>
            <CardDescription>Top Performing Businesses</CardDescription>
          </CardHeader>
          <CardContent>
            <ChartContainer config={{}} className="w-full max-h-[600px]">
              <BarChart
                accessibilityLayer
                data={[
                  ...(topBusinesses ?? []).map((topBusiness) => ({
                    x: topBusiness.name,
                    y: topBusiness.weight,
                    color: 'var(--color-primary)',
                    fill: 'var(--color-primary)',
                  })),
                ]}
                margin={{ top: 20, right: 20, left: 20, bottom: 20 }}
                dataKey="y"
              >
                <CartesianGrid />
                <XAxis
                  dataKey="x"
                  tickLine={false}
                  tickMargin={10}
                  axisLine={false}
                />
                <YAxis
                  dataKey="y"
                  tickLine={false}
                  tickMargin={10}
                  axisLine={false}
                />

                <Bar isAnimationActive={false} dataKey="y" radius={8}>
                  <LabelList
                    dataKey="y"
                    position="top"
                    fill="var(--foreground)"
                    fontWeight="semibold"
                    fontFamily="montserrat"
                    fontSize="12pt"
                    formatter={(value: string) => `${value}kg`}
                  />
                </Bar>

                <ChartLegend
                  content={
                    <ChartLegendContent
                      className="grid grid-cols-6 gap-x-3 gap-y-1"
                      verticalAlign="bottom"
                      payload={[...(topBusinesses ?? [])].map(
                        (item, index) => ({
                          id: `${index}`,
                          value: item.name,
                          color: 'var(--color-primary)',
                          dataKey: 'x',
                        })
                      )}
                    />
                  }
                />
              </BarChart>
            </ChartContainer>
          </CardContent>
        </Card>

        <Card className="bg-background">
          <CardHeader>
            <CardDescription>Material Types</CardDescription>
          </CardHeader>
          <CardContent>
            <ChartContainer config={{}} className="w-full max-h-[600px]">
              <PieChart accessibilityLayer>
                <Pie
                  data={[...(materialTypes ?? [])]}
                  dataKey="weight"
                  nameKey="type"
                  fill="var(--primary)"
                  label={renderCustomizedLabel}
                  isAnimationActive={false}
                />
              </PieChart>
            </ChartContainer>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
