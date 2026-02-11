/// <reference types="@cloudflare/workers-types" />

interface Env {
  BUCKET: R2Bucket;
  IMAGE_BASE_URL: string;
  IMAGE_FOLDER: string;
}

const IMAGE_EXTENSIONS = /\.(jpg|jpeg|png|gif|webp|avif)$/i;

export const onRequestGet: PagesFunction<Env> = async (context) => {
  const { BUCKET, IMAGE_BASE_URL, IMAGE_FOLDER } = context.env;

  const listed = await BUCKET.list({
    prefix: `${IMAGE_FOLDER}/`,
    limit: 1000,
  });

  const images = listed.objects
    .filter((obj) => IMAGE_EXTENSIONS.test(obj.key))
    .sort((a, b) => b.uploaded.getTime() - a.uploaded.getTime())
    .map((obj) => ({
      key: obj.key,
      size: obj.size,
      uploaded: obj.uploaded.toISOString(),
    }));

  return Response.json({
    images,
    baseURL: IMAGE_BASE_URL,
    folder: IMAGE_FOLDER,
  });
};
