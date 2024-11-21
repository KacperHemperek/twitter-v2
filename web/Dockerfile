FROM node:20 AS build 

ARG VITE_API_URL

ENV VITE_API_URL=$VITE_API_URL

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

WORKDIR /web

COPY . .

RUN pnpm install

RUN pnpm build

FROM nginx:alpine AS deploy

COPY --from=build /web/dist /usr/share/nginx/html
COPY docker.nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
