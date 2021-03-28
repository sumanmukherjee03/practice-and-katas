class HealthCheckController < ApplicationController
  def hz
    render json: { ok: true , node: "It's alive!"}
  end
end
